package ht

import (
	"transformer/src/lib/dao/schema"
)

func (processor HashTableProcessor[T]) ProcessGetRequest(request IRequest[T]) IResult[T] {

	table, bucket, err := processor.hashManager.Locate(request.ObjectId())

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &HashTableGetResult[T]{object: nil, err: err}
	}
	return &HashTableGetResult[T]{object: &bucket.object, err: err}
}

type HashTableGetRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	objectId uint64
}

func (request *HashTableGetRequest[T]) Type() RequestType {

	return GetRequest
}

func (request *HashTableGetRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *HashTableGetRequest[T]) ObjectId() uint64 {

	return request.objectId
}

func (request *HashTableGetRequest[T]) Object() *T {

	return nil
}

type HashTableGetResult[T schema.FixedSizeIdentifiable] struct {
	object *T

	err error
}

func (result *HashTableGetResult[T]) Object() *T {

	return result.object
}

func (result *HashTableGetResult[T]) Error() error {

	return result.err
}
