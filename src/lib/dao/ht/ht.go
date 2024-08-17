package ht

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

type Hash struct {
	tableGroup int

	tableNumber int

	position int
}

type ITSFHashTable[T schema.FixedSizeIdentifiable] interface {
	Save(object T) error

	Get(id uint64) (*T, error)

	Update(object T) error

	Delete(id uint64) error
}

type TSFHashTable[T schema.FixedSizeIdentifiable] struct {
	id uint64

	fileManager fm.IFileManager

	tableManager IHashTableManager[T]

	bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader]

	objectIO daoio.IFixedSizeDaoIO[T]
}

func New[T schema.FixedSizeIdentifiable](path, descriptor string, id uint64) TSFHashTable[T] {

	fileManager := fm.NewFileContainer(path, descriptor)

	tableManager := NewTableManager[T](fileManager)

	bucketHeaderIO := daoio.FixedSizeDaoIO[HashTableBucketHeader]{}

	objectIO := daoio.FixedSizeDaoIO[T]{}

	return TSFHashTable[T]{id, fileManager, tableManager, bucketHeaderIO, objectIO}
}

func Inject[T schema.FixedSizeIdentifiable](id uint64, fileManager fm.IFileManager, tableManager IHashTableManager[T], bucketHeaderIO daoio.IFixedSizeDaoIO[HashTableBucketHeader], objectIO daoio.IFixedSizeDaoIO[T]) TSFHashTable[T] {

	return TSFHashTable[T]{id, fileManager, tableManager, bucketHeaderIO, objectIO}
}
