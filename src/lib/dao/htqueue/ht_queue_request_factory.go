package htqueue

import (
	"transformer/src/lib/dao/queue"
	"transformer/src/lib/dao/schema"
)

type HashTableRequestFactory[T schema.FixedSizeIdentifiable] struct {
}

func (factory HashTableRequestFactory[T]) CreateSaveRequest(requestId uint64, object T) queue.IRequest[T] {

	return &HashTableSaveRequest[T]{requestId, object}
}

func (factory HashTableRequestFactory[T]) CreateGetRequest(requestId uint64, objectId uint64) queue.IRequest[T] {

	return &HashTableGetRequest[T]{requestId, objectId}
}

func (factory HashTableRequestFactory[T]) CreateUpdateRequest(requestId uint64, object T) queue.IRequest[T] {

	return &HashTableUpdateRequest[T]{requestId, object}
}

func (factory HashTableRequestFactory[T]) CreateDeleteRequest(requestId uint64, objectId uint64) queue.IRequest[T] {

	return &HashTableDeleteRequest[T]{requestId, objectId}
}
