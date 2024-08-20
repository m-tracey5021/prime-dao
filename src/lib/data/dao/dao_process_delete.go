package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

// func (dao TSFDao[T]) Delete(id uint64) (int, error) {

// 	index, err := dao.indexHashTable.Get(id)

// 	if err != nil {

// 		return 0, err
// 	}
// 	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

// 	defer dao.fileManager.CloseAndUnlock(table, &err)

// 	if err != nil {

// 		return 0, err
// 	}
// 	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

// 		return 0, err
// 	}
// 	deletedSize, err := dao.objectIO.Delete(table)

// 	if err != nil {

// 		return 0, err
// 	}
// 	dao.tableManager.AlterCache(id, RemoveObject)

// 	dao.UpdateMetadataForDeletion(table, index.fileId)

// 	dao.UpdateIndexes(table, index.fileId, index.filePosition)

// 	// could run these two in parallel
// 	if err := dao.tableManager.SaveManagingInfo(); err != nil {

// 		return 0, err
// 	}
// 	if err := dao.indexHashTable.Delete(index.id); err != nil {

// 		return 0, err
// 	}
// 	return deletedSize, err
// }

type DeleteProcessor[T schema.Identifiable] struct {
	id uint64

	fileManager fm.IFileManager

	metadataManager *DaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

type ProcessDeleteResult[T schema.Identifiable] struct {
	sizeDeleted int

	err error
}

func (result *ProcessDeleteResult[T]) Object() *T {

	return nil
}

func (result *ProcessDeleteResult[T]) Error() error {

	return result.err
}

func (processor DeleteProcessor[T]) Process() queue.IResult[T] {

	index, err := processor.metadataManager.GetIndex(processor.id)

	if err != nil {

		return &ProcessDeleteResult[T]{0, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &ProcessDeleteResult[T]{0, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &ProcessDeleteResult[T]{0, err}
	}
	sizeDeleted, err := processor.objectIO.Delete(table)

	if err != nil {

		return &ProcessDeleteResult[T]{0, err}
	}
	if err := processor.metadataManager.UpdateForDeletion(table, *index); err != nil {

		return &ProcessDeleteResult[T]{0, err}
	}
	if err := processor.metadataManager.DeleteIndex(index.id); err != nil {

		return &ProcessDeleteResult[T]{0, err}
	}
	return &ProcessDeleteResult[T]{sizeDeleted, err}
}
