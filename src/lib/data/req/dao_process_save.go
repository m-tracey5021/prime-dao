package req

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoSaveRequest[T schema.Identifiable] struct {
	requestId uint64

	object T
}

func (request *DaoSaveRequest[T]) Type() queue.RequestType {

	return queue.SaveRequest
}

func (request *DaoSaveRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *DaoSaveRequest[T]) ObjectId() uint64 {

	return request.object.Id()
}

func (request *DaoSaveRequest[T]) Object() *T {

	return &request.object
}

type DaoSaveResult[T any] struct {
	object *T

	sizeWritten int

	err error
}

func (result *DaoSaveResult[T]) Object() *T {

	return nil
}

func (result *DaoSaveResult[T]) Error() error {

	return result.err
}

// func (processor DaoRequestProcessor[T]) ProcessSaveRequest(request queue.IRequest[T]) queue.IResult[T] {

// 	object, sizeWritten, err := processor.dao.Save(*request.Object())

// 	return &DaoSaveResult[T]{object, sizeWritten, err}
// }
