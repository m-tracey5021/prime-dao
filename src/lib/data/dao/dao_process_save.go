package dao

import (
	"sync"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type SaveRequest[T schema.Identifiable] struct {
	requestId uint64

	object T

	dependencies chan queue.RequestContext

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor SaveRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor SaveRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor SaveRequest[T]) Dependencies() chan queue.RequestContext {

	return processor.dependencies
}

func (processor SaveRequest[T]) Process(wg *sync.WaitGroup, context *queue.QueueContext[T]) queue.IResult[T] {

	defer wg.Done()

	// context.Complete(processor.requestId)

	// see 'Get' for wait

	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, processor.metadataManager.AvailableTable())

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
	processor.metadataManager.UpdateForSave(table, object, uint64(position))

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
