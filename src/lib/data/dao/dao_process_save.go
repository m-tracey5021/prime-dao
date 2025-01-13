package dao

import (
	"fmt"
	"prime-dao/src/lib/data/dataio"
	"prime-dao/src/lib/data/fm"
	"prime-dao/src/lib/data/queue"
	"prime-dao/src/lib/data/schema"
	"sync"
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

	availableTable := processor.metadataManager.UpdateMetadataPreSave()

	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, availableTable)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &SaveResult[T]{0, err}
	}
	position, err := processor.fileManager.Size(table)

	if err := processor.fileManager.GoTo(position, table); err != nil {

		return &SaveResult[T]{0, err}
	}
	size, err := processor.objectIO.WriteSizePrefixed(table, processor.object)

	if err != nil {

		return &SaveResult[T]{0, err}
	}
	processor.metadataManager.UpdateMetadataPostSave(availableTable, processor.object, uint64(position))

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
