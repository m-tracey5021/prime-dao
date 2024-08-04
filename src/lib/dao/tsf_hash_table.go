package dao

import (
	"errors"
	. "transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
	"unsafe"
)

type ITSFHashTable[T FixedSizeIdentifiable] interface {
	Save(object T) error

	Get(id uint64) (*T, error)

	Update(object T) error

	Delete(id uint64) error
}

type TSFHashTable[T FixedSizeIdentifiable] struct {
	id uint64

	bucketSize int

	tableSize int

	// Controls file operations and maintains file name consistency
	fileContainer IFileManager

	identifierCache HashTableIdentifierCache

	// Controls IO of the cache that keeps track of object ids
	cacheIO daoio.IDaoIO[HashTableIdentifierCache]

	// Controls IO of the bucket headers which contain info on hashing collisions
	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader]

	// Controls IO of the actual object to be hashed and saved in the main table
	objectIO daoio.IFixedSizeDaoIO[T]
}

func Default[T FixedSizeIdentifiable](
	id uint64,

	bucketSize int,

	tableSize int,

	fileContainer IFileManager,

	identifierCache HashTableIdentifierCache,

	cacheIO daoio.IDaoIO[HashTableIdentifierCache],

	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader],

	objectIO daoio.IFixedSizeDaoIO[T],

) TSFHashTable[T] {

	return TSFHashTable[T]{id, bucketSize, tableSize, fileContainer, identifierCache, cacheIO, bucketHeaderIO, objectIO}
}

func New[T FixedSizeIdentifiable](
	id uint64,

	fileManager IFileManager,

	cacheIO daoio.IDaoIO[HashTableIdentifierCache],

	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader],

	objectIO daoio.IFixedSizeDaoIO[T],

) (ITSFHashTable[T], error) {

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	managingFile, err := fileManager.Open(HashTableManagingFile, id)

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
	return &TSFHashTable[T]{id, bucketSize, tableSize, fileManager, *identifierCache, cacheIO, bucketHeaderIO, objectIO}, err
}

func Wrapped[T FixedSizeIdentifiable](fileManager IFileManager, id uint64) (ITSFHashTable[T], error) {

	// fileManager := NewFileContainer(path, descriptor)

	cacheIO := daoio.DaoIO[HashTableIdentifierCache]{}

	bucketHeaderIO := daoio.FixedSizeDaoIO[HashTableBucketHeader]{}

	objectIO := daoio.FixedSizeDaoIO[T]{}

	return New(id, fileManager, cacheIO, bucketHeaderIO, objectIO)
}

func (dao *TSFHashTable[T]) hash(id uint64) int {

	hash := (int(id)*2 + 1) % dao.tableSize

	return hash * dao.bucketSize
}

func (dao *TSFHashTable[T]) NewCollisionTableId() uint64 {

	return smallestMissing(dao.identifierCache.CollisionTableIds)
}
