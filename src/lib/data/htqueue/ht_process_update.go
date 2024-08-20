package htqueue

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type HashTableUpdateRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	object T
}

func (request *HashTableUpdateRequest[T]) Type() queue.RequestType {

	return queue.UpdateRequest
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

func (processor HashTableRequestProcessor[T]) ProcessUpdateRequest(request queue.IRequest[T]) queue.IResult[T] {

	err := processor.hashTable.Update(*request.Object())

	return &HashTableUpdateResult[T]{err}
}
