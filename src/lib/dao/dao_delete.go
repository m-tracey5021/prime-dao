package dao

import (
	"errors"
	"transformer/src/lib/dao/fm"
)

func (dao TSFDao[T]) Delete(object DaoIdentifiable[T]) (int, error) {

	index, err := dao.indexHashTable.Get(object.Id)

	if err != nil {

		return 0, err
	}
	table, err := dao.fileManager.Open(fm.DaoTable, index.fileId)

	if err != nil {

		return 0, err
	}
	defer func() {

		if innerErr := dao.fileManager.Close(table); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return 0, err
	}
	deletedSize, err := dao.objectIO.Delete(table)

	if err != nil {

		return 0, err
	}
	if err := dao.SaveCacheAlteration(object.Id, dao.RemoveObjectIdFromCache); err != nil {

		return 0, err
	}
	dao.UpdateMetadataForDeletion(table, index.fileId)

	dao.UpdateIndexes(table, index.fileId, index.filePosition)

	if err := dao.indexHashTable.Delete(index.id); err != nil {

		return 0, err
	}
	return deletedSize, err
}
