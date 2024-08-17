package dao

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/daoio"
)

func (dao TSFDao[T]) Get(id uint64) (*T, error) {

	index, err := dao.indexHashTable.Get(id)

	if err != nil {

		return nil, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.id)

	defer dao.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return nil, err
	}
	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return nil, err
	}
	object, err := daoio.ReadSizePrefixed[T](table)

	if err != nil {

		return object, err
	}
	return nil, err
}
