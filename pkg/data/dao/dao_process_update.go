package dao

import (
	"fmt"
	"sync"

	"github.com/m-tracey5021/prime-dao/pkg/data/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type UpdateRequest[T schema.Identifiable] struct {
	requestId uint64

	object T

	dependencies chan queue.QueueDependency

	numberOfDependencies int

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor *UpdateRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *UpdateRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor *UpdateRequest[T]) Dependencies() chan queue.QueueDependency {

	return processor.dependencies
}

func (processor *UpdateRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *UpdateRequest[T]) ProcessWithDependencies(wg *sync.WaitGroup, context *queue.QueueDependencyResolver[T]) queue.IResult[T] {

	defer wg.Done()

	var depWaitGroup sync.WaitGroup

	depWaitGroup.Add(processor.numberOfDependencies)

	go func() {

		for reqCtx := range processor.dependencies {

			reqCtx.RequestWaitGroup.Wait()

			depWaitGroup.Done()

			fmt.Printf("request %v waited for dependency %v", processor.requestId, reqCtx.RequestId)
		}
	}()
	depWaitGroup.Wait()

	close(processor.dependencies)

	return processor.Process()
}

func (processor *UpdateRequest[T]) Process() queue.IResult[T] {

	index, err := processor.metadataManager.GetIndex(processor.object.Id())

	if err != nil {

		return &UpdateResult[T]{0, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &UpdateResult[T]{0, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &UpdateResult[T]{0, err}
	}
	updatedSize, err := processor.objectIO.Update(table, processor.object)

	if err != nil {

		return &UpdateResult[T]{0, err}
	}
	updatedFilePosition := int(index.filePosition) + updatedSize

	if err := processor.metadataManager.UpdateIndexes(table, index.fileId, uint64(updatedFilePosition)); err != nil {

		return &UpdateResult[T]{0, err}
	}
	processor.metadataManager.UpdateCache(processor.object)

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
