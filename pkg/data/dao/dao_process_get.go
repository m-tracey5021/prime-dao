package dao

import (
	"fmt"
	"prime-dao/pkg/data/dataio"
	"prime-dao/pkg/data/fm"
	"prime-dao/pkg/data/queue"
	"prime-dao/pkg/data/schema"
	"sync"
)

type GetRequest[T schema.Identifiable] struct {
	requestId uint64

	objectId uint64

	dependencies chan queue.QueueDependency

	numberOfDependencies int

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor *GetRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *GetRequest[T]) ObjectId() uint64 {

	return processor.objectId
}

func (processor *GetRequest[T]) Dependencies() chan queue.QueueDependency {

	return processor.dependencies
}

func (processor *GetRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *GetRequest[T]) ProcessWithDependencies(wg *sync.WaitGroup, context *queue.QueueDependencyResolver[T]) queue.IResult[T] {

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

func (processor *GetRequest[T]) Process() queue.IResult[T] {

	cached := processor.metadataManager.GetCached(processor.objectId)

	if cached != nil {

		return &GetResult[T]{cached, nil}
	}
	index, err := processor.metadataManager.GetIndex(processor.objectId)

	if err != nil {

		return &GetResult[T]{nil, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &GetResult[T]{nil, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &GetResult[T]{nil, err}
	}
	object, err := processor.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return &GetResult[T]{nil, err}
	}
	processor.metadataManager.SaveToCache(object)

	return &GetResult[T]{&object, err}
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
