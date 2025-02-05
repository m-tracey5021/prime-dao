package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type DaoDeleteRequest[T schema.Identifiable] struct {
	requestId uint64

	objectId uuid.UUID

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoDeleteRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoDeleteRequest[T]) ObjectId() uuid.UUID {

	return processor.objectId
}

func (processor *DaoDeleteRequest[T]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoDeleteRequest[T]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoDeleteRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoDeleteRequest[T]) Process(dao *Dao[T]) IResult[T] {

	sizeDeleted, err := dao.Delete(processor.objectId)

	return &DeleteResult[T]{sizeDeleted, err}
}

type DeleteResult[T schema.Identifiable] struct {
	sizeDeleted int

	err error
}

func (result *DeleteResult[T]) Object() *T {

	return nil
}

func (result *DeleteResult[T]) Size() int {

	return result.sizeDeleted
}

func (result *DeleteResult[T]) Error() error {

	return result.err
}
