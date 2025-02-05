package dao

import "github.com/m-tracey5021/prime-dao/pkg/fm"

func (dao *Dao[T]) Get(objectId uint64) (*T, error) {

	cached := dao.metadata.Cache.Get(objectId, &dao.cacheMutex)

	if cached != nil {

		return cached, nil
	}
	index, err := dao.indexHashTable.Get(objectId)

	if err != nil {

		return nil, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer dao.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return nil, err
	}
	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return nil, err
	}
	object, err := dao.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return nil, err
	}
	dao.metadata.Cache.Save(object, &dao.cacheMutex)

	return &object, err
}
