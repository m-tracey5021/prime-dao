package dao

import "github.com/m-tracey5021/prime-dao/pkg/schema"

type IProcessableRequest[T schema.Identifiable] interface {
	RequestId() uint64

	ObjectId() uint64

	Dependencies() chan uint64

	SetNumDeps(int)

	GetNumDeps() int

	Process(*Dao[T]) IResult[T]
}

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}
