package htqueue

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type HashTableGetRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	objectId uint64
}

func (request *HashTableGetRequest[T]) Type() queue.RequestType {

	return queue.GetRequest
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

func (processor HashTableRequestProcessor[T]) ProcessGetRequest(request queue.IRequest[T]) queue.IResult[T] {

	object, err := processor.hashTable.Get(request.ObjectId())

	return &HashTableGetResult[T]{object, err}
}
