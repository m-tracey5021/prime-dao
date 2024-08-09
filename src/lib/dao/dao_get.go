package dao

import (
	"transformer/src/lib/dao/fm"
)

func (dao TSFDao[T]) Get(id uint64) (*DaoIdentifiable[T], error) {

	index, err := dao.indexHashTable.Get(id)

	if err != nil {

		return nil, err
	}
	table, err := dao.fileManager.Open(fm.DaoTable, index.id)

	if err != nil {

		return nil, err
	}
	defer dao.fileManager.Close(table, &err)

	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return nil, err
	}
	object, err := dao.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return nil, err
	}
	return object, err
}
