package ht

import (
	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type Hash struct {
	tableGroup int

	tableNumber int

	position int
}

type ITSFHashTable[T schema.FixedSizeIdentifiable] interface {
	Save(object T) error

	Get(id uint64) (*T, error)

	GetSome(ids ...uint64) ([]*T, error)

	Update(object T) error

	Delete(id uint64) error
}

type TSFHashTable[T schema.FixedSizeIdentifiable] struct {
	fileManager fm.IFileManager

	tableManager IHashTableManager[T]

	bucketHeaderIO dataio.IFixedSizeDataIO[HashTableBucketHeader]

	objectIO dataio.IFixedSizeDataIO[T]
}

func New[T schema.FixedSizeIdentifiable](path, descriptor string, id uint64) TSFHashTable[T] {

	fileManager := fm.NewFileManager(path, descriptor)

	tableManager := NewTableManager[T](fileManager)

	bucketHeaderIO := dataio.FixedSizeDataIO[HashTableBucketHeader]{}

	objectIO := dataio.FixedSizeDataIO[T]{}

	return TSFHashTable[T]{fileManager, tableManager, bucketHeaderIO, objectIO}
}

func FromFileManager[T schema.FixedSizeIdentifiable](fileManager fm.IFileManager) ITSFHashTable[T] {

	tableManager := NewTableManager[T](fileManager)

	bucketHeaderIO := dataio.FixedSizeDataIO[HashTableBucketHeader]{}

	objectIO := dataio.FixedSizeDataIO[T]{}

	return &TSFHashTable[T]{fileManager, tableManager, bucketHeaderIO, objectIO}
}
