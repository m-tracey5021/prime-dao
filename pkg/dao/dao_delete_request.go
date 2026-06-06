package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type DaoDeleteRequest[T schema.Orderable[U], U constraints.Ordered] struct {
	requestId uint64

	objectId uuid.UUID

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoDeleteRequest[T, U]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoDeleteRequest[T, U]) ObjectId() uuid.UUID {

	return processor.objectId
}

func (processor *DaoDeleteRequest[T, U]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoDeleteRequest[T, U]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoDeleteRequest[T, U]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoDeleteRequest[T, U]) Process(dao *Dao[T, U]) IResult[T] {

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
