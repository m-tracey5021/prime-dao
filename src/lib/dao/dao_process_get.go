package dao

import (
	"transformer/src/lib/dao/queue"
)

type DaoGetRequest[T any] struct {
	requestId uint64

	objectId uint64
}

func (request *DaoGetRequest[T]) Type() queue.RequestType {

	return queue.GetRequest
}

func (request *DaoGetRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *DaoGetRequest[T]) ObjectId() uint64 {

	return request.objectId
}

func (request *DaoGetRequest[T]) Object() *T {

	return nil
}

type DaoGetResult[T any] struct {
	object *T

	err error
}

func (result *DaoGetResult[T]) Object() *T {

	return result.object
}

func (result *DaoGetResult[T]) Error() error {

	return result.err
}

func (processor DaoRequestProcessor[T]) ProcessGetRequest(request queue.IRequest[T]) queue.IResult[T] {

	object, err := processor.dao.Get(request.ObjectId())

	return &DaoGetResult[T]{object, err}
}
