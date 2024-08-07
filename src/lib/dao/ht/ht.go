package ht

import (
	"errors"
	"fmt"
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

	// Controls IO of the table headers which contain info on individual collision tables
	tableHeaderIO daoio.IFixedSizeDaoIO[HashTableCollisionTableHeader]

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

	tableHeaderIO daoio.IFixedSizeDaoIO[HashTableCollisionTableHeader],

	objectIO daoio.IFixedSizeDaoIO[T],

) TSFHashTable[T] {

	return TSFHashTable[T]{id, bucketSize, tableSize, maxCollisions, fileContainer, identifierCache, cacheIO, bucketHeaderIO, tableHeaderIO, objectIO}
}

func Initialise[T schema.FixedSizeIdentifiable](
	id uint64,

	fileManager fm.IFileManager,

	cacheIO daoio.IDaoIO[HashTableIdentifierCache],

	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader],

	tableHeaderIO daoio.IFixedSizeDaoIO[HashTableCollisionTableHeader],

	objectIO daoio.IFixedSizeDaoIO[T],

) (ITSFHashTable[T], error) {

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	managingFile, err := fileManager.Open(fm.HashTableManagingFile, id)

	if err != nil {

		return &TSFHashTable[T]{}, err
	}
	defer func() {

		if innerErr := fileManager.Close(managingFile); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

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
	return &TSFHashTable[T]{id, bucketSize, tableSize, maxCollisions, fileManager, *identifierCache, cacheIO, bucketHeaderIO, tableHeaderIO, objectIO}, err
}

func New[T schema.FixedSizeIdentifiable](fileManager fm.IFileManager, id uint64) (ITSFHashTable[T], error) {

	cacheIO := daoio.DaoIO[HashTableIdentifierCache]{}

	bucketHeaderIO := daoio.FixedSizeDaoIO[HashTableBucketHeader]{}

	tableHeaderIO := daoio.FixedSizeDaoIO[HashTableCollisionTableHeader]{}

	objectIO := daoio.FixedSizeDaoIO[T]{}

	return Initialise(id, fileManager, cacheIO, bucketHeaderIO, tableHeaderIO, objectIO)
}

func (ht *TSFHashTable[T]) hash(id uint64) int {

	hash := (int(id)*2 + 1) % ht.tableSize

	return hash * ht.bucketSize
}

func (ht *TSFHashTable[T]) hashCollision(id uint64, hash int) (int, string) {

	// Calculate the order of the collision
	collisionPosition := (int(id) / ht.tableSize) * ht.bucketSize

	// Calculate the file/partition in which the collision will be stored
	collisionTable := collisionPosition / ht.maxCollisions

	collisionTableId := fmt.Sprintf("%v_%v", hash, collisionTable)

	return collisionPosition, collisionTableId
}

func (ht *TSFHashTable[T]) NewCollisionTableId() uint64 {

	return lib.NewId(ht.identifierCache.CollisionTableIds)
}
