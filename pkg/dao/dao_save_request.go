package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type DaoSaveRequest[T schema.Orderable[U], U constraints.Ordered] struct {
	requestId uint64

	object T

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoSaveRequest[T, U]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoSaveRequest[T, U]) ObjectId() uuid.UUID {

	return processor.object.Id()
}

func (processor *DaoSaveRequest[T, U]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoSaveRequest[T, U]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoSaveRequest[T, U]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoSaveRequest[T, U]) Process(dao *Dao[T, U]) IResult[T] {

	size, err := dao.Save(processor.object)

	return &SaveResult[T]{size, err}
}

type SaveResult[T schema.Identifiable] struct {
	sizeWritten int

	err error
}

func (result *SaveResult[T]) Object() *T {

	return nil
}

func (result *SaveResult[T]) Size() int {

	return result.sizeWritten
}

func (result *SaveResult[T]) Error() error {

	return result.err
}
