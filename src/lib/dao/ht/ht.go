package ht

import (
	"errors"
	"io"
	"os"
	"transformer/src/lib"
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
	"unsafe"
)

type ITSFHashTable[T schema.FixedSizeIdentifiable] interface {
	Save(object T) error

	Get(id uint64) (*T, error)

	Update(object T) error

	Delete(id uint64) error
}

type TSFHashTable[T schema.FixedSizeIdentifiable] struct {
	id uint64

	bucketSize int

	tableSize int

	maxCollisions int

	// Controls file operations and maintains file name consistency
	fileContainer fm.IFileManager

	identifierCache HashTableIdentifierCache

	// Controls IO of the cache that keeps track of object ids
	cacheIO daoio.IDaoIO[HashTableIdentifierCache]

	// Controls IO of the bucket headers which contain info on hashing collisions
	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader]

	// Controls IO of the actual object to be hashed and saved in the main table
	objectIO daoio.IFixedSizeDaoIO[T]
}

func Default[T schema.FixedSizeIdentifiable](
	id uint64,

	bucketSize int,

	tableSize int,

	maxCollisions int,

	fileContainer fm.IFileManager,

	identifierCache HashTableIdentifierCache,

	cacheIO daoio.IDaoIO[HashTableIdentifierCache],

	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader],

	objectIO daoio.IFixedSizeDaoIO[T],

) TSFHashTable[T] {

	return TSFHashTable[T]{id, bucketSize, tableSize, maxCollisions, fileContainer, identifierCache, cacheIO, bucketHeaderIO, objectIO}
}

func Initialise[T schema.FixedSizeIdentifiable](
	id uint64,

	fileManager fm.IFileManager,

	cacheIO daoio.IDaoIO[HashTableIdentifierCache],

	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader],

	objectIO daoio.IFixedSizeDaoIO[T],

) (ITSFHashTable[T], error) {

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	managingFile, err := fileManager.Open(fm.HashTableManagingFile, id)

	defer fileManager.Close(managingFile, &err)

	if err != nil {

		return &TSFHashTable[T]{}, err
	}
	size, err := fileManager.Size(managingFile)

	if err != nil {

		return &TSFHashTable[T]{}, err
	}
	var identifierCache *HashTableIdentifierCache

	if size > 0 {

		if identifierCache, err = cacheIO.ReadSizePrefixed(managingFile); err != nil {

			return &TSFHashTable[T]{}, err
		}

	} else {

		collisionTableIds := make([]uint64, 0)

		identifierCache = &HashTableIdentifierCache{collisionTableIds}

		if _, err := cacheIO.WriteSizePrefixed(managingFile, *identifierCache); err != nil {

			return &TSFHashTable[T]{}, err
		}
	}
	return &TSFHashTable[T]{id, bucketSize, tableSize, maxCollisions, fileManager, *identifierCache, cacheIO, bucketHeaderIO, objectIO}, err
}

func New[T schema.FixedSizeIdentifiable](fileManager fm.IFileManager, id uint64) (ITSFHashTable[T], error) {

	cacheIO := daoio.DaoIO[HashTableIdentifierCache]{}

	bucketHeaderIO := daoio.FixedSizeDaoIO[HashTableBucketHeader]{}

	objectIO := daoio.FixedSizeDaoIO[T]{}

	return Initialise(id, fileManager, cacheIO, bucketHeaderIO, objectIO)
}

func (ht *TSFHashTable[T]) hash(id uint64) int {

	hash := int(id) % ht.tableSize

	return hash * ht.bucketSize
}

func (ht *TSFHashTable[T]) hashCollision(id uint64) Hash {

	// Calculate the table group to store the data in
	hash := int(id) % ht.tableSize

	// Calculate the order of the hash i.e. where it sits in relation to the others if hashed
	order := int(id) / ht.tableSize

	// Calculate the file/partition in which the data is stored
	tableNumber := order / ht.maxCollisions

	// Calculate the actual position in the file based on bucket size and table number
	position := (order - (ht.maxCollisions * tableNumber)) * ht.bucketSize

	return Hash{hash, tableNumber, position}
}

func (ht *TSFHashTable[T]) NewCollisionTableId() uint64 {

	return lib.NewId(ht.identifierCache.CollisionTableIds)
}

func (ht *TSFHashTable[T]) ReadBucketHeader(position int, table *os.File) (HashTableBucketHeader, error) {

	if err := ht.fileContainer.GoTo(position, table); err != nil {

		return HashTableBucketHeader{}, err
	}
	bucketHeader, err := ht.bucketHeaderIO.Read(table)

	if err != nil {

		if errors.Is(err, io.EOF) {

			return HashTableBucketHeader{}, nil
		}
	}
	return bucketHeader, err
}

func (ht *TSFHashTable[T]) LocateForTableAndPosition(id uint64, table *os.File, position int) (*HashTableBucket[T], error) {

	bucketHeader, err := ht.ReadBucketHeader(position, table)

	if err != nil {

		return nil, err
	}
	if bucketHeader.occupied {

		objectPosition, err := ht.fileContainer.CurrentPosition(table)

		if err != nil {

			return nil, err
		}
		object, err := ht.objectIO.Read(table)

		if err != nil {

			return nil, err
		}
		if id == object.Id() {

			return &HashTableBucket[T]{position, objectPosition, bucketHeader, object}, nil

		} else {

			return nil, nil
		}

	} else {

		return nil, nil
	}
}

func (ht *TSFHashTable[T]) LocateEmptyForTableAndPosition(id uint64, table *os.File, position int) (*int, error) {

	bucketHeader, err := ht.ReadBucketHeader(position, table)

	if err != nil {

		return nil, err
	}
	if bucketHeader.occupied {

		object, err := ht.objectIO.Read(table)

		if err != nil {

			return nil, err
		}
		if id == object.Id() {

			return nil, ObjectAlreadyExists

		} else {

			return nil, nil
		}

	} else {

		return &position, nil
	}
}

func (ht *TSFHashTable[T]) Locate(id uint64) (*os.File, HashTableBucket[T], error) {

	closeTable := true

	hash := ht.hashCollision(id)

	table, err := ht.fileContainer.Open(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

	defer ht.CloseConditionally(&closeTable, table, &err)

	if err != nil {

		return nil, HashTableBucket[T]{}, err
	}
	location, err := ht.LocateForTableAndPosition(id, table, hash.position)

	if err != nil {

		return nil, HashTableBucket[T]{}, err
	}
	if location != nil {

		closeTable = false

		return table, *location, err
	}
	return nil, HashTableBucket[T]{}, ObjectDoesNotExist
}

func (ht *TSFHashTable[T]) LocateEmpty(id uint64) (*os.File, int, error) {

	closeTable := true

	hash := ht.hashCollision(id)

	table, err := ht.fileContainer.Open(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

	defer ht.CloseConditionally(&closeTable, table, &err)

	if err != nil {

		return nil, 0, err
	}
	location, err := ht.LocateEmptyForTableAndPosition(id, table, hash.position)

	if err != nil {

		return nil, 0, err
	}
	if location != nil {

		closeTable = false

		return table, *location, nil
	}
	return nil, 0, BucketOccupied
}

// Closes the file conditionally so that defer will run only if it needs to
func (ht *TSFHashTable[T]) CloseConditionally(close *bool, table *os.File, err *error) error {

	if *close {

		return ht.fileContainer.Close(table, err)
	}
	return *err
}
