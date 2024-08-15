package ht

import (
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

func (processor HashTableProcessor[T]) ProcessSaveRequest(request IRequest[T]) IResult[T] {

	table, emptyBucket, err := processor.hashManager.LocateEmpty(request.ObjectId())

	defer processor.fileManager.CloseAndUnlock(table, &err)

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

type HashTableSaveRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	object T
}

func (request *HashTableSaveRequest[T]) Type() RequestType {

	return SaveRequest
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
	err error
}

func (result *HashTableSaveResult[T]) Object() *T {

	return nil
}

func (result *HashTableSaveResult[T]) Error() error {

	return result.err
}
