package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DeleteRequest[T schema.Identifiable] struct {
	objectId uint64

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor DeleteRequest[T]) ObjectId() uint64 {

	return processor.objectId
}

func (processor DeleteRequest[T]) Process() queue.IResult[T] {

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
