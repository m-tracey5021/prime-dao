package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type DaoGetRequest[T schema.Identifiable] struct {
	requestId uint64

	objectId uuid.UUID

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoGetRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoGetRequest[T]) ObjectId() uuid.UUID {

	return processor.objectId
}

func (processor *DaoGetRequest[T]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoGetRequest[T]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoGetRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoGetRequest[T]) Process(dao *Dao[T]) IResult[T] {

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
