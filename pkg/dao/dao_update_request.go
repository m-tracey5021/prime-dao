package dao

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type DaoUpdateRequest[T schema.Orderable[U], U constraints.Ordered] struct {
	requestId uint64

	object T

	dependencies chan uint64

	numberOfDependencies int
}

func (processor *DaoUpdateRequest[T, U]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DaoUpdateRequest[T, U]) ObjectId() uuid.UUID {

	return processor.object.Id()
}

func (processor *DaoUpdateRequest[T, U]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DaoUpdateRequest[T, U]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DaoUpdateRequest[T, U]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DaoUpdateRequest[T, U]) Process(dao *Dao[T, U]) IResult[T] {

	updatedSize, err := dao.Update(processor.object)

	return &UpdateResult[T]{updatedSize, err}
}

type UpdateResult[T schema.Identifiable] struct {
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
