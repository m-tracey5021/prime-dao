package htqueue

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type HashTableSaveRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	object T
}

func (request *HashTableSaveRequest[T]) Type() queue.RequestType {

	return queue.SaveRequest
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

func (processor HashTableRequestProcessor[T]) ProcessSaveRequest(request queue.IRequest[T]) queue.IResult[T] {

	err := processor.hashTable.Save(*request.Object())

	return &HashTableSaveResult[T]{err}
}
