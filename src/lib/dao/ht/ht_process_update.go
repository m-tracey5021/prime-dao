package ht

import (
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

func (processor HashTableProcessor[T]) ProcessUpdateRequest(request IRequest[T]) IResult[T] {

	table, bucket, err := processor.hashManager.Locate(request.ObjectId())

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &HashTableUpdateResult[T]{err: err}
	}
	if err := processor.fileManager.GoTo(bucket.objectLocation, table); err != nil {

		return &HashTableUpdateResult[T]{err: err}
	}
	if err := daoio.Write(table, *request.Object()); err != nil {

		return &HashTableUpdateResult[T]{err: err}
	}
	return &HashTableUpdateResult[T]{err: err}
}

type HashTableUpdateRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	object T
}

func (request *HashTableUpdateRequest[T]) Type() RequestType {

	return UpdateRequest
}

func (request *HashTableUpdateRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *HashTableUpdateRequest[T]) ObjectId() uint64 {

	return request.object.Id()
}

func (request *HashTableUpdateRequest[T]) Object() *T {

	return &request.object
}

type HashTableUpdateResult[T schema.FixedSizeIdentifiable] struct {
	err error
}

func (result *HashTableUpdateResult[T]) Object() *T {

	return nil
}

func (result *HashTableUpdateResult[T]) Error() error {

	return result.err
}
