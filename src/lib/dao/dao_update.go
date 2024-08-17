package dao

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/daoio"
)

func (dao TSFDao[T]) Update(object T) (int, error) {

	index, err := dao.indexHashTable.Get(object.Id())

	if err != nil {

		return 0, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.id)

	defer dao.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return 0, err
	}
	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return 0, err
	}
	updatedSize, err := daoio.UpdateSizePrefixed(table, object)

	if err != nil {

		return 0, err
	}
	updatedFilePosition := int(index.filePosition) + updatedSize

	if err := dao.UpdateIndexes(table, index.fileId, uint64(updatedFilePosition)); err != nil {

		return 0, err
	}
	return updatedSize, err
}
