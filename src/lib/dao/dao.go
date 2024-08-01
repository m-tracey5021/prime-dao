package dao

import (
	"errors"
	. "transformer/src/lib/dao/schema"
	. "transformer/src/lib/dao_io"
	"unsafe"
)

// This FixedSizeDao only works for objects that are a constant size, i.e. no lists of maps etc.
type FixedSizeDao[T FixedSizeIdentifiable] struct {
	bucketSize int

	tableSize int

	// Controls file operations and maintains file name consistency
	fileContainer IDaoFileContainer

	identifierCache DaoIdentifierCache

	// Controls IO of the cache that keeps track of object ids
	cacheIO IDaoIO[DaoIdentifierCache]

	// Controls IO of the bucket headers which contain info on hashing collisions
	bucketHeaderIO IFixedSizeDaoIO[DaoBucketHeader]

	// Controls IO of the actual object to be hashed and saved in the main table
	objectIO IFixedSizeDaoIO[T]
}

func Default[T FixedSizeIdentifiable](
	bucketSize int,

	tableSize int,

	fileContainer IDaoFileContainer,

	identifierCache DaoIdentifierCache,

	cacheIO IDaoIO[DaoIdentifierCache],

	bucketHeaderIO IFixedSizeDaoIO[DaoBucketHeader],

	objectIO IFixedSizeDaoIO[T],

) FixedSizeDao[T] {

	return FixedSizeDao[T]{bucketSize, tableSize, fileContainer, identifierCache, cacheIO, bucketHeaderIO, objectIO}
}

func New[T FixedSizeIdentifiable](

	fileContainer IDaoFileContainer,

	cacheIO IDaoIO[DaoIdentifierCache],

	bucketHeaderIO IFixedSizeDaoIO[DaoBucketHeader],

	objectIO IFixedSizeDaoIO[T],

) (FixedSizeDao[T], error) {

	bucketSize := int(unsafe.Sizeof(*new(DaoBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	managingFile, err := fileContainer.ManagingFile()

	defer func() {

		if innerErr := fileContainer.Close(managingFile); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	if err != nil {

		return FixedSizeDao[T]{}, err
	}
	size, err := fileContainer.Size(managingFile)

	if err != nil {

		return FixedSizeDao[T]{}, err
	}
	var identifierCache *DaoIdentifierCache

	if size > 0 {

		if identifierCache, err = cacheIO.ReadSizePrefixed(managingFile); err != nil {

			return FixedSizeDao[T]{}, err
		}

	} else {

		objectIds := make([]uint64, 0)

		fileIds := make([]uint64, 0)

		identifierCache = &DaoIdentifierCache{objectIds, fileIds}

		if _, err := cacheIO.WriteSizePrefixed(managingFile, *identifierCache); err != nil {

			return FixedSizeDao[T]{}, err
		}
	}
	return FixedSizeDao[T]{bucketSize, tableSize, fileContainer, *identifierCache, cacheIO, bucketHeaderIO, objectIO}, nil
}

func Wrapped[T FixedSizeIdentifiable](path, descriptor string) (FixedSizeDao[T], error) {

	fileContainer := DaoFileContainer{path, descriptor}

	cacheIO := DaoIO[DaoIdentifierCache]{}

	bucketHeaderIO := FixedSizeDaoIO[DaoBucketHeader]{}

	objectIO := FixedSizeDaoIO[T]{}

	return New(fileContainer, cacheIO, bucketHeaderIO, objectIO)
}

func (dao *FixedSizeDao[T]) hash(id uint64) int {

	hash := (int(id)*2 + 1) % dao.tableSize

	return hash * dao.bucketSize
}

func (dao *FixedSizeDao[T]) NewCollisionTableId() uint64 {

	return smallestMissing(dao.identifierCache.FileIds)
}

func (dao *FixedSizeDao[T]) NewObjectId() uint64 {

	return smallestMissing(dao.identifierCache.ObjectIds)
}
