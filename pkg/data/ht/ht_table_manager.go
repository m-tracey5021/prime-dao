package ht

import (
	"errors"
	"io"
	"os"
	"unsafe"

	"github.com/m-tracey5021/prime-dao/pkg/data/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type IHashTableManager[T schema.FixedSizeIdentifiable] interface {
	ComputeHash(id uint64) Hash

	Locate(id uint64) (*os.File, HashTableBucket[T], error)

	LocateEmpty(id uint64) (*os.File, int, error)
}

type HashTableManager[T schema.FixedSizeIdentifiable] struct {
	bucketSize int

	tableSize int

	maxCollisions int

	fileManager fm.IFileManager

	bucketHeaderIO dataio.IFixedSizeDataIO[HashTableBucketHeader]

	objectIO dataio.IFixedSizeDataIO[T]
}

func NewTableManager[T schema.FixedSizeIdentifiable](fileManager fm.IFileManager) IHashTableManager[T] {

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	bucketHeaderIO := dataio.FixedSizeDataIO[HashTableBucketHeader]{}

	objectIO := dataio.FixedSizeDataIO[T]{}

	return &HashTableManager[T]{bucketSize, tableSize, maxCollisions, fileManager, bucketHeaderIO, objectIO}
}

func InjectTableManager[T schema.FixedSizeIdentifiable](fileManager fm.IFileManager, bucketHeaderIO dataio.IFixedSizeDataIO[HashTableBucketHeader], objectIO dataio.IFixedSizeDataIO[T]) HashTableManager[T] {

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	return HashTableManager[T]{bucketSize, tableSize, maxCollisions, fileManager, bucketHeaderIO, objectIO}
}

func (manager HashTableManager[T]) ComputeHash(id uint64) Hash {

	// Calculate the table group to store the data in
	hash := int(id) % manager.tableSize

	// Calculate the order of the hash i.e. where it sits in relation to the others if hashed
	order := int(id) / manager.tableSize

	// Calculate the file/partition in which the data is stored
	tableNumber := order / manager.maxCollisions

	// Calculate the actual position in the file based on bucket size and table number
	position := (order - (manager.maxCollisions * tableNumber)) * manager.bucketSize

	return Hash{hash, tableNumber, position}
}

func (manager *HashTableManager[T]) ReadBucketHeader(position int, table *os.File) (HashTableBucketHeader, error) {

	if err := manager.fileManager.GoTo(position, table); err != nil {

		return HashTableBucketHeader{}, err
	}
	bucketHeader, err := manager.bucketHeaderIO.Read(table)

	if err != nil {

		if errors.Is(err, io.EOF) {

			return HashTableBucketHeader{}, nil
		}
	}
	return bucketHeader, err
}

func (manager *HashTableManager[T]) LocateForTableAndPosition(id uint64, table *os.File, position int) (*HashTableBucket[T], error) {

	bucketHeader, err := manager.ReadBucketHeader(position, table)

	if err != nil {

		return nil, err
	}
	if bucketHeader.occupied {

		objectPosition, err := manager.fileManager.CurrentPosition(table)

		if err != nil {

			return nil, err
		}
		object, err := manager.objectIO.Read(table)

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

func (manager *HashTableManager[T]) LocateEmptyForTableAndPosition(id uint64, table *os.File, position int) (*int, error) {

	bucketHeader, err := manager.ReadBucketHeader(position, table)

	if err != nil {

		return nil, err
	}
	if bucketHeader.occupied {

		object, err := manager.objectIO.Read(table)

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

func (manager *HashTableManager[T]) Locate(id uint64) (*os.File, HashTableBucket[T], error) {

	closeTable := true

	hash := manager.ComputeHash(id)

	table, err := manager.fileManager.OpenAndLock(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

	defer manager.CloseConditionally(&closeTable, table, &err)

	if err != nil {

		return nil, HashTableBucket[T]{}, err
	}
	location, err := manager.LocateForTableAndPosition(id, table, hash.position)

	if err != nil {

		return nil, HashTableBucket[T]{}, err
	}
	if location != nil {

		closeTable = false

		return table, *location, err
	}
	return nil, HashTableBucket[T]{}, ObjectDoesNotExist
}

func (manager *HashTableManager[T]) LocateEmpty(id uint64) (*os.File, int, error) {

	closeTable := true

	hash := manager.ComputeHash(id)

	table, err := manager.fileManager.OpenAndLock(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

	defer manager.CloseConditionally(&closeTable, table, &err)

	if err != nil {

		return nil, 0, err
	}
	location, err := manager.LocateEmptyForTableAndPosition(id, table, hash.position)

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
func (manager *HashTableManager[T]) CloseConditionally(close *bool, table *os.File, err *error) error {

	if *close {

		return manager.fileManager.CloseAndUnlock(table, err)
	}
	return *err
}
