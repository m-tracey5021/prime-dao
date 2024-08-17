package dao

import (
	"transformer/src/lib/dao/queue"
	"transformer/src/lib/dao/schema"
)

type DaoUpdateRequest[T schema.Identifiable] struct {
	requestId uint64

	object T
}

func (request *DaoUpdateRequest[T]) Type() queue.RequestType {

	return queue.UpdateRequest
}

func (request *DaoUpdateRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *DaoUpdateRequest[T]) ObjectId() uint64 {

	return request.object.Id()
}

func (request *DaoUpdateRequest[T]) Object() *T {

	return &request.object
}

type DaoUpdateResult[T any] struct {
	updatedSize int

	err error
}

func (result *DaoUpdateResult[T]) Object() *T {

	return nil
}

func (result *DaoUpdateResult[T]) Error() error {

	return result.err
}

func (processor DaoRequestProcessor[T]) ProcessUpdateRequest(request queue.IRequest[T]) queue.IResult[T] {

	updatedSize, err := processor.dao.Update(*request.Object())

	return &DaoUpdateResult[T]{updatedSize, err}
}
