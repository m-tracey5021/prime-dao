package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

// func (dao TSFDao[T]) Update(object T) (int, error) {

// 	index, err := dao.indexHashTable.Get(object.Id())

// 	if err != nil {

// 		return 0, err
// 	}
// 	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.id)

// 	defer dao.fileManager.CloseAndUnlock(table, &err)

// 	if err != nil {

// 		return 0, err
// 	}
// 	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

// 		return 0, err
// 	}
// 	updatedSize, err := dao.objectIO.Update(table, object)

// 	if err != nil {

// 		return 0, err
// 	}
// 	updatedFilePosition := int(index.filePosition) + updatedSize

// 	if err := dao.UpdateIndexes(table, index.fileId, uint64(updatedFilePosition)); err != nil {

// 		return 0, err
// 	}
// 	return updatedSize, err
// }

type UpdateProcessor[T schema.Identifiable] struct {
	object T

	fileManager fm.IFileManager

	metadataManager *DaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

type ProcessUpdateResult[T schema.Identifiable] struct {
	sizeUpdated int

	err error
}

func (result *ProcessUpdateResult[T]) Object() *T {

	return nil
}

func (result *ProcessUpdateResult[T]) Error() error {

	return result.err
}

func (processor UpdateProcessor[T]) Process() queue.IResult[T] {

	index, err := processor.metadataManager.GetIndex(processor.object.Id())

	if err != nil {

		return &ProcessUpdateResult[T]{0, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.id)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &ProcessUpdateResult[T]{0, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &ProcessUpdateResult[T]{0, err}
	}
	updatedSize, err := processor.objectIO.Update(table, processor.object)

	if err != nil {

		return &ProcessUpdateResult[T]{0, err}
	}
	updatedFilePosition := int(index.filePosition) + updatedSize

	if err := processor.metadataManager.UpdateIndexes(table, index.fileId, uint64(updatedFilePosition)); err != nil {

		return &ProcessUpdateResult[T]{0, err}
	}
	return &ProcessUpdateResult[T]{updatedSize, err}
}
