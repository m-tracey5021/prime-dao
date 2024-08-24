package dao

import (
	"sync"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DeleteRequest[T schema.Identifiable] struct {
	requestId uint64

	objectId uint64

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor DeleteRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor DeleteRequest[T]) ObjectId() uint64 {

	return processor.objectId
}

func (processor DeleteRequest[T]) Process(wg *sync.WaitGroup, context *queue.QueueContext[T]) queue.IResult[T] {

	defer wg.Done()

	// context.Complete(processor.requestId)

	// see 'Get' for wait

	index, err := processor.metadataManager.GetIndex(processor.objectId)

	if err != nil {

		return &DeleteResult[T]{0, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &DeleteResult[T]{0, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &DeleteResult[T]{0, err}
	}
	sizeDeleted, err := processor.objectIO.Delete(table)

	if err != nil {

		return &DeleteResult[T]{0, err}
	}
	if err := processor.metadataManager.UpdateForDeletion(table, *index); err != nil {

		return &DeleteResult[T]{0, err}
	}
	if err := processor.metadataManager.DeleteIndex(index.id); err != nil {

		return &DeleteResult[T]{0, err}
	}
	return &DeleteResult[T]{sizeDeleted, err}
}

type DeleteResult[T schema.Identifiable] struct {
	sizeDeleted int

	err error
}

func (result *DeleteResult[T]) Object() *T {

	return nil
}

func (result *DeleteResult[T]) Size() int {

	return result.sizeDeleted
}

func (result *DeleteResult[T]) Error() error {

	return result.err
}
