package dao

import (
	"sync"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type UpdateRequest[T schema.Identifiable] struct {
	requestId uint64

	object T

	dependencies chan queue.RequestContext

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor UpdateRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor UpdateRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor UpdateRequest[T]) Dependencies() chan queue.RequestContext {

	return processor.dependencies
}

func (processor UpdateRequest[T]) Process(wg *sync.WaitGroup, context *queue.QueueContext[T]) queue.IResult[T] {

	defer wg.Done()

	// context.Complete(processor.requestId)

	// see 'Get' for wait

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
