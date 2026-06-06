package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type IProcessableRequest[T schema.Orderable[U], U constraints.Ordered] interface {
	RequestId() uint64

	ObjectId() uuid.UUID

	Dependencies() chan uint64

	SetNumDeps(int)

	GetNumDeps() int

	Process(*Dao[T, U]) IResult[T]
}

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}
