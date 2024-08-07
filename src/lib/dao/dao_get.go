package dao

import (
	"errors"
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

		return nil, err
	}
	object, err := dao.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return nil, err
	}
	return object, err
}
