package req

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoRequestFactory[T schema.Identifiable] struct {
}

func (factory DaoRequestFactory[T]) CreateSaveRequest(requestId uint64, object T) queue.IRequest[T] {

	return &DaoSaveRequest[T]{requestId, object}
}

func (factory DaoRequestFactory[T]) CreateGetRequest(requestId uint64, objectId uint64) queue.IRequest[T] {

	return &DaoGetRequest[T]{requestId, objectId}
}

func (factory DaoRequestFactory[T]) CreateUpdateRequest(requestId uint64, object T) queue.IRequest[T] {

	return &DaoUpdateRequest[T]{requestId, object}
}

func (factory DaoRequestFactory[T]) CreateDeleteRequest(requestId uint64, objectId uint64) queue.IRequest[T] {

	return &DaoDeleteRequest[T]{requestId, objectId}
}
