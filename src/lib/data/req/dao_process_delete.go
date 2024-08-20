package req

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoDeleteRequest[T schema.Identifiable] struct {
	requestId uint64

	objectId uint64
}

func (request *DaoDeleteRequest[T]) Type() queue.RequestType {

	return queue.DeleteRequest
}

func (request *DaoDeleteRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *DaoDeleteRequest[T]) ObjectId() uint64 {

	return request.objectId
}

func (request *DaoDeleteRequest[T]) Object() *T {

	return nil
}

type DaoDeleteResult[T any] struct {
	deletedSize int

	err error
}

func (result *DaoDeleteResult[T]) Object() *T {

	return nil
}

func (result *DaoDeleteResult[T]) Error() error {

	return result.err
}

// func (processor DaoRequestProcessor[T]) ProcessDeleteRequest(request queue.IRequest[T]) queue.IResult[T] {

// 	deletedSize, err := processor.dao.Delete(request.ObjectId())

// 	return &DaoDeleteResult[T]{deletedSize, err}
// }
