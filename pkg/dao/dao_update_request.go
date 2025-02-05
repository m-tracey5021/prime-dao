package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type DaoUpdateRequest[T schema.DescribedIdentifiable] struct {
	requestId uint64

	object T

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoUpdateRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoUpdateRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor *DaoUpdateRequest[T]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoUpdateRequest[T]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoUpdateRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoUpdateRequest[T]) Process(dao *Dao[T]) IResult[T] {

	updatedSize, err := dao.Update(processor.object)

	return &UpdateResult[T]{updatedSize, err}
}

type UpdateResult[T schema.DescribedIdentifiable] struct {
	sizeUpdated int

	err error
}

func (result *UpdateResult[T]) Object() *T {

	return nil
}

func (result *UpdateResult[T]) Size() int {

	return result.sizeUpdated
}

func (result *UpdateResult[T]) Error() error {

	return result.err
}
