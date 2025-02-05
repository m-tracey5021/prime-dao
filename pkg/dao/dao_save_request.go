package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type DaoSaveRequest[T schema.DescribedIdentifiable] struct {
	requestId uint64

	object T

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoSaveRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoSaveRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor *DaoSaveRequest[T]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoSaveRequest[T]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoSaveRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoSaveRequest[T]) Process(dao *Dao[T]) IResult[T] {

	size, err := dao.Save(processor.object)

	return &SaveResult[T]{size, err}
}

type SaveResult[T schema.DescribedIdentifiable] struct {
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
