package ht

import (
	"errors"
	"io"
	"os"
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

type BucketManager[T schema.FixedSizeIdentifiable] struct {
	fileManager fm.IFileManager

	hashMetrics HashMetrics
}

func (manager *BucketManager[T]) ReadBucketHeader(position int, table *os.File) (HashTableBucketHeader, error) {

	if err := manager.fileManager.GoTo(position, table); err != nil {

		return HashTableBucketHeader{}, err
	}
	bucketHeader, err := daoio.Read[HashTableBucketHeader](table)

	if err != nil {

		if errors.Is(err, io.EOF) {

			return HashTableBucketHeader{}, nil
		}
	}
	return bucketHeader, err
}

func (manager *BucketManager[T]) LocateForTableAndPosition(id uint64, table *os.File, position int) (*HashTableBucket[T], error) {

	bucketHeader, err := manager.ReadBucketHeader(position, table)

	if err != nil {

		return nil, err
	}
	if bucketHeader.occupied {

		objectPosition, err := manager.fileManager.CurrentPosition(table)

		if err != nil {

			return nil, err
		}
		object, err := daoio.Read[T](table)

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

func (manager *BucketManager[T]) LocateEmptyForTableAndPosition(id uint64, table *os.File, position int) (*int, error) {

	bucketHeader, err := manager.ReadBucketHeader(position, table)

	if err != nil {

		return nil, err
	}
	if bucketHeader.occupied {

		object, err := daoio.Read[T](table)

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

func (manager *BucketManager[T]) Locate(id uint64) (*os.File, HashTableBucket[T], error) {

	closeTable := true

	hash := manager.hashMetrics.ComputeHash(id)

	table, err := manager.fileManager.Open(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

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

func (manager *BucketManager[T]) LocateEmpty(id uint64) (*os.File, int, error) {

	closeTable := true

	hash := manager.hashMetrics.ComputeHash(id)

	table, err := manager.fileManager.Open(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

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
func (manager *BucketManager[T]) CloseConditionally(close *bool, table *os.File, err *error) error {

	if *close {

		return manager.fileManager.Close(table, err)
	}
	return *err
}
