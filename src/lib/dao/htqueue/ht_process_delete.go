package htqueue

import (
	"transformer/src/lib/dao/queue"
	"transformer/src/lib/dao/schema"
)

type HashTableDeleteRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	objectId uint64
}

func (request *HashTableDeleteRequest[T]) Type() queue.RequestType {

	return queue.DeleteRequest
}

func (request *HashTableDeleteRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *HashTableDeleteRequest[T]) ObjectId() uint64 {

	return request.objectId
}

func (request *HashTableDeleteRequest[T]) Object() *T {

	return nil
}

type HashTableDeleteResult[T schema.FixedSizeIdentifiable] struct {
	err error
}

func (result *HashTableDeleteResult[T]) Object() *T {

	return nil
}

func (result *HashTableDeleteResult[T]) Error() error {

	return result.err
}

func (processor HashTableRequestProcessor[T]) ProcessDeleteRequest(request queue.IRequest[T]) queue.IResult[T] {

	err := processor.hashTable.Delete(request.ObjectId())

	return &HashTableDeleteResult[T]{err}
}
