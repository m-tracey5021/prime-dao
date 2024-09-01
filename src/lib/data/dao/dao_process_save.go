package dao

import (
	"fmt"
	"sync"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type SaveRequest[T schema.Identifiable] struct {
	requestId uint64

	object T

	dependencies chan queue.QueueDependency

	numberOfDependencies int

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor *SaveRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *SaveRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor *SaveRequest[T]) Dependencies() chan queue.QueueDependency {

	return processor.dependencies
}

func (processor *SaveRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *SaveRequest[T]) ProcessWithDependencies(wg *sync.WaitGroup, context *queue.QueueDependencyResolver[T]) queue.IResult[T] {

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

func (processor *SaveRequest[T]) Process() queue.IResult[T] {

	availableTable := processor.metadataManager.AvailableTable()

	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, availableTable)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &SaveResult[T]{nil, 0, err}
	}
	object := processor.metadataManager.AssignId(processor.object)

	position, err := processor.fileManager.Size(table)

	if err := processor.fileManager.GoTo(position, table); err != nil {

		return &SaveResult[T]{nil, 0, err}
	}
	size, err := processor.objectIO.WriteSizePrefixed(table, object)

	if err != nil {

		return &SaveResult[T]{nil, 0, err}
	}
	processor.metadataManager.UpdateForSave(availableTable, object, uint64(position))

	return &SaveResult[T]{&object, size, err}
}

type SaveResult[T schema.Identifiable] struct {
	object *T

	sizeWritten int

	err error
}

func (result *SaveResult[T]) Object() *T {

	return result.object
}

func (result *SaveResult[T]) Size() int {

	return result.sizeWritten
}

func (result *SaveResult[T]) Error() error {

	return result.err
}
