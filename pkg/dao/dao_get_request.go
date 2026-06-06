package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type DaoGetRequest[T schema.Orderable[U], U constraints.Ordered] struct {
	requestId uint64

	objectId uuid.UUID

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoGetRequest[T, U]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoGetRequest[T, U]) ObjectId() uuid.UUID {

	return processor.objectId
}

func (processor *DaoGetRequest[T, U]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoGetRequest[T, U]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoGetRequest[T, U]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoGetRequest[T, U]) Process(dao *Dao[T, U]) IResult[T] {

	object, err := dao.Get(processor.objectId)

	return &GetResult[T]{object, err}
}

type GetResult[T schema.Identifiable] struct {
	object *T

	err error
}

func (result *GetResult[T]) Object() *T {

	return result.object
}

func (result *GetResult[T]) Size() int {

	return 0
}

func (result *GetResult[T]) Error() error {

	return result.err
}
