package dao

import (
	"errors"
	"transformer/src/lib/dao/fm"
)

func (dao TSFDao[T]) Save(object T) (int, error) {

	table, err := dao.fileManager.Open(fm.DaoTable, dao.managingInfo.AvailableTable)

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

	objectId := dao.NewObjectId()

	identifiable := DaoIdentifiable[T]{objectId, object}

	position, err := dao.fileManager.Size(table)

	if err := dao.fileManager.GoTo(position, table); err != nil {

		return 0, err
	}
	size, err := dao.objectIO.WriteSizePrefixed(table, identifiable)

	if err != nil {

		return 0, err
	}
	if err := dao.SaveCacheAlteration(objectId, dao.AddObjectIdToCache); err != nil {

		return 0, err
	}
	index := DaoIndex{objectId, dao.managingInfo.AvailableTable, uint64(position)}

	if err := dao.indexHashTable.Save(index); err != nil {

		return 0, err
	}
	if err := dao.UpdateMetadataForSave(dao.managingInfo.AvailableTable); err != nil {

		return 0, err
	}
	return size, err
}
