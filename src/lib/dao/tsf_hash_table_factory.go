package dao

import (
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

type ITSFHashTableFactory[T schema.FixedSizeIdentifiable] interface {
	Create(fileManager IFileManager, id uint64) (ITSFHashTable[T], error)
}

type TSFHashTableFactory[T schema.FixedSizeIdentifiable] struct{}

func (hashTableFactory TSFHashTableFactory[T]) Create(fileManager IFileManager, id uint64) (ITSFHashTable[T], error) {

	// fileManager := NewFileContainer(path, descriptor)

	cacheIO := daoio.DaoIO[HashTableIdentifierCache]{}

	bucketHeaderIO := daoio.FixedSizeDaoIO[HashTableBucketHeader]{}

	objectIO := daoio.FixedSizeDaoIO[T]{}

	return New(id, fileManager, cacheIO, bucketHeaderIO, objectIO)
}
