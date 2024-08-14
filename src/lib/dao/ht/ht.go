package ht

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
	"unsafe"
)

type ITSFHashTable[T schema.FixedSizeIdentifiable] interface {
	Save(object T) error

	// SaveConcurrent(objects ...T) error

	Get(id uint64) (*T, error)

	// GetConcurrent(ids ...uint64) ([]*T, error)

	Update(object T) error

	Delete(id uint64) error
}

type TSFHashTable[T schema.FixedSizeIdentifiable] struct {
	id uint64

	hashMetrics HashMetrics

	saveQueue *CommandQueue[T]
}

func Default[T schema.FixedSizeIdentifiable](
	id uint64,

	bucketSize int,

	tableSize int,

	maxCollisions int,

	fileContainer fm.IFileManager,

	identifierCache HashTableIdentifierCache,

	cacheIO daoio.IDaoIO[HashTableIdentifierCache],

	saveQueue *CommandQueue[T],

) TSFHashTable[T] {

	hashMetrics := HashMetrics{bucketSize, tableSize, maxCollisions}

	return TSFHashTable[T]{id, hashMetrics, saveQueue}
}

func Initialise[T schema.FixedSizeIdentifiable](
	id uint64,

	fileManager fm.IFileManager,

	saveQueue *CommandQueue[T],

) (ITSFHashTable[T], error) { // TODO get rid of error

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	hashMetrics := HashMetrics{bucketSize, tableSize, maxCollisions}

	return &TSFHashTable[T]{id, hashMetrics, saveQueue}, nil
}

func New[T schema.FixedSizeIdentifiable](fileManager fm.IFileManager, id uint64) (ITSFHashTable[T], error) {

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	hashMetrics := HashMetrics{bucketSize, tableSize, maxCollisions}

	bucketLocator := BucketManager[T]{fileManager, hashMetrics}

	saveProcessor := SaveProcessor[T]{fileManager, bucketLocator}

	saveQueue := NewCommandQueue(10, 5, saveProcessor)

	return &TSFHashTable[T]{id, hashMetrics, saveQueue}, nil
}

// func (ht *TSFHashTable[T]) ReadBucketHeader(position int, table *os.File) (HashTableBucketHeader, error) {

// 	if err := ht.fileManager.GoTo(position, table); err != nil {

// 		return HashTableBucketHeader{}, err
// 	}
// 	bucketHeader, err := ht.bucketHeaderIO.Read(table)

// 	if err != nil {

// 		if errors.Is(err, io.EOF) {

// 			return HashTableBucketHeader{}, nil
// 		}
// 	}
// 	return bucketHeader, err
// }

// func (ht *TSFHashTable[T]) LocateForTableAndPosition(id uint64, table *os.File, position int) (*HashTableBucket[T], error) {

// 	bucketHeader, err := ht.ReadBucketHeader(position, table)

// 	if err != nil {

// 		return nil, err
// 	}
// 	if bucketHeader.occupied {

// 		objectPosition, err := ht.fileManager.CurrentPosition(table)

// 		if err != nil {

// 			return nil, err
// 		}
// 		object, err := ht.objectIO.Read(table)

// 		if err != nil {

// 			return nil, err
// 		}
// 		if id == object.Id() {

// 			return &HashTableBucket[T]{position, objectPosition, bucketHeader, object}, nil

// 		} else {

// 			return nil, nil
// 		}

// 	} else {

// 		return nil, nil
// 	}
// }

// func (ht *TSFHashTable[T]) LocateEmptyForTableAndPosition(id uint64, table *os.File, position int) (*int, error) {

// 	bucketHeader, err := ht.ReadBucketHeader(position, table)

// 	if err != nil {

// 		return nil, err
// 	}
// 	if bucketHeader.occupied {

// 		object, err := ht.objectIO.Read(table)

// 		if err != nil {

// 			return nil, err
// 		}
// 		if id == object.Id() {

// 			return nil, ObjectAlreadyExists

// 		} else {

// 			return nil, nil
// 		}

// 	} else {

// 		return &position, nil
// 	}
// }

// func (ht *TSFHashTable[T]) Locate(id uint64) (*os.File, HashTableBucket[T], error) {

// 	closeTable := true

// 	hash := ht.hashCollision(id)

// 	table, err := ht.fileManager.Open(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

// 	defer ht.CloseConditionally(&closeTable, table, &err)

// 	if err != nil {

// 		return nil, HashTableBucket[T]{}, err
// 	}
// 	location, err := ht.LocateForTableAndPosition(id, table, hash.position)

// 	if err != nil {

// 		return nil, HashTableBucket[T]{}, err
// 	}
// 	if location != nil {

// 		closeTable = false

// 		return table, *location, err
// 	}
// 	return nil, HashTableBucket[T]{}, ObjectDoesNotExist
// }

// func (ht *TSFHashTable[T]) LocateEmpty(id uint64) (*os.File, int, error) {

// 	closeTable := true

// 	hash := ht.hashCollision(id)

// 	table, err := ht.fileManager.Open(fm.HashTableCollisionTable, uint64(hash.tableGroup), uint64(hash.tableNumber))

// 	defer ht.CloseConditionally(&closeTable, table, &err)

// 	if err != nil {

// 		return nil, 0, err
// 	}
// 	location, err := ht.LocateEmptyForTableAndPosition(id, table, hash.position)

// 	if err != nil {

// 		return nil, 0, err
// 	}
// 	if location != nil {

// 		closeTable = false

// 		return table, *location, nil
// 	}
// 	return nil, 0, BucketOccupied
// }

// // Closes the file conditionally so that defer will run only if it needs to
// func (ht *TSFHashTable[T]) CloseConditionally(close *bool, table *os.File, err *error) error {

// 	if *close {

// 		return ht.fileManager.Close(table, err)
// 	}
// 	return *err
// }
