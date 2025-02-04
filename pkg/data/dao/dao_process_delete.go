package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/data/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type DeleteRequest[T schema.Identifiable] struct {
	requestId uint64

	objectId uint64

	dependencies chan uint64

	numberOfDependencies int

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor *DeleteRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *DeleteRequest[T]) ObjectId() uint64 {

	return processor.objectId
}

func (processor *DeleteRequest[T]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *DeleteRequest[T]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *DeleteRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *DeleteRequest[T]) Process() queue.IResult[T] {

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
	if err := processor.metadataManager.UpdateMetadataForDeletion(table, *index); err != nil {

		return &DeleteResult[T]{0, err}
	}
	if err := processor.metadataManager.DeleteIndex(index.id); err != nil {

		return &DeleteResult[T]{0, err}
	}
	processor.metadataManager.DeleteFromCache(processor.objectId)

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
