package dao

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/daoio"
)

func (dao TSFDao[T]) Delete(id uint64) (int, error) {

	index, err := dao.indexHashTable.Get(id)

	if err != nil {

		return 0, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer dao.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return 0, err
	}
	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return 0, err
	}
	deletedSize, err := daoio.DeleteSizePrefixed[T](table)

	if err != nil {

		return 0, err
	}
	if err := dao.SaveCacheAlteration(id, dao.RemoveObjectIdFromCache); err != nil {

		return 0, err
	}
	dao.UpdateMetadataForDeletion(table, index.fileId)

	dao.UpdateIndexes(table, index.fileId, index.filePosition)

	if err := dao.indexHashTable.Delete(index.id); err != nil {

		return 0, err
	}
	return deletedSize, err
}
