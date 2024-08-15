package ht

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
)

type RequestType int

const (
	SaveRequest = iota

	GetRequest

	UpdateRequest

	DeleteRequest
)

type IRequest[T any] interface {
	Type() RequestType

	RequestId() uint64

	ObjectId() uint64

	Object() *T
}

type IResult[T any] interface {
	Object() *T

	Error() error
}

type IRequestProcessor[T any] interface {
	Process(IRequest[T]) IResult[T]
}

type HashTableProcessor[T schema.FixedSizeIdentifiable] struct {
	fileManager fm.IFileManager

	hashManager HashManager[T]
}

func (processor HashTableProcessor[T]) Process(request IRequest[T]) IResult[T] {

	switch request.Type() {

	case SaveRequest:

		return processor.ProcessSaveRequest(request)

	case GetRequest:

		return processor.ProcessGetRequest(request)

	case UpdateRequest:

		return processor.ProcessUpdateRequest(request)

	case DeleteRequest:

		return processor.ProcessDeleteRequest(request)

	default:

		return nil
	}
}
