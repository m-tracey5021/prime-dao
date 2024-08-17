package dao

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/daoio"
)

func (dao TSFDao[T]) Save(object T) (*T, int, error) {

	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, dao.managingInfo.AvailableTable)

	defer dao.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return nil, 0, err
	}
	object = dao.AssignId(object)

	position, err := dao.fileManager.Size(table)

	if err := dao.fileManager.GoTo(position, table); err != nil {

		return nil, 0, err
	}
	size, err := daoio.WriteSizePrefixed(table, object)

	if err != nil {

		return nil, 0, err
	}
	if err := dao.SaveCacheAlteration(object.Id(), dao.AddObjectIdToCache); err != nil {

		return nil, 0, err
	}
	index := DaoIndex{object.Id(), dao.managingInfo.AvailableTable, uint64(position)}

	if err := dao.indexHashTable.Save(index); err != nil {

		return nil, 0, err
	}
	if err := dao.UpdateMetadataForSave(dao.managingInfo.AvailableTable); err != nil {

		return nil, 0, err
	}
	return &object, size, err
}
