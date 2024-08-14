package ht

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

type SaveProcessor[T schema.FixedSizeIdentifiable] struct {
	fileManager fm.IFileManager

	bucketLocator BucketManager[T]
}

func (processor SaveProcessor[T]) Process(request IRequest[T]) IResult[T] {

	table, emptyBucket, err := processor.bucketLocator.LocateEmpty(request.ObjectId())

	defer processor.fileManager.Close(table, &err)

	if err != nil {

		return &HashTableSaveResult[T]{err: err}
	}
	if err := processor.fileManager.GoTo(emptyBucket, table); err != nil {

		return &HashTableSaveResult[T]{err: err}
	}
	bucketHeader := HashTableBucketHeader{true, false, 0}

	if err := daoio.Write(table, bucketHeader); err != nil {

		return &HashTableSaveResult[T]{err: err}
	}
	if err := daoio.Write(table, *request.Object()); err != nil {

		return &HashTableSaveResult[T]{err: err}
	}
	return &HashTableSaveResult[T]{err: err}
}

// func (processor SaveProcessor[T]) MarshallRequest(request IRequest[T]) *HashTableSaveRequest[T] {

// 	if concreteResult, ok := request.(*HashTableSaveRequest[T]); ok {

// 		return concreteResult

// 	} else {

// 		panic("type assertion failed for HashTableSaveRequest")
// 	}
// }

// func (processor SaveProcessor[T]) MarshallResult(result IResult[T]) *HashTableSaveResult[T] {

// 	if concreteResult, ok := result.(*HashTableSaveResult[T]); ok {

// 		return concreteResult

// 	} else {

// 		panic("type assertion failed for HashTableSaveResult")
// 	}
// }

type HashTableSaveRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	object T
}

func (request *HashTableSaveRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *HashTableSaveRequest[T]) ObjectId() uint64 {

	return request.object.Id()
}

func (request *HashTableSaveRequest[T]) Object() *T {

	return &request.object
}

type HashTableSaveResult[T schema.FixedSizeIdentifiable] struct {
	associatedRequestId uint64

	err error
}

func (result *HashTableSaveResult[T]) AssociatedRequestId() uint64 {

	return result.associatedRequestId
}

func (result *HashTableSaveResult[T]) Object() *T {

	return nil
}

func (result *HashTableSaveResult[T]) Error() error {

	return result.err
}

func (result *HashTableSaveResult[T]) SetAssociatedRequestId(id uint64) {

	result.associatedRequestId = id
}
