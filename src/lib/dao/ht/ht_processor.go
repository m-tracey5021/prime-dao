package ht

import "transformer/src/lib/dao/schema"

type IRequest[T schema.FixedSizeIdentifiable] interface {
	RequestId() uint64

	ObjectId() uint64

	Object() *T
}

type IResult[T schema.FixedSizeIdentifiable] interface {
	AssociatedRequestId() uint64

	Object() *T

	Error() error

	SetAssociatedRequestId(uint64)
}

type IRequestProcessor[T schema.FixedSizeIdentifiable] interface {
	Process(IRequest[T]) IResult[T]
}
