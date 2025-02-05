package dao

import (
	"errors"
	"io"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
)

func (dao *Dao[T]) Get(objectId uuid.UUID) (*T, error) {

	cached := dao.metadata.Cache.Get(objectId, &dao.cacheMutex)

	if cached != nil {

		return cached, nil
	}
	index, err := dao.indexHashTable.Get(objectId)

	if err != nil {

		return nil, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoObjectFile, index.fileId)

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

func (dao *Dao[T]) GetAll() ([]T, error) {

	objects := []T{}

	var err error

	for _, objectFileId := range dao.metadata.TableIdStore.AllIds() {

		table, err := dao.fileManager.OpenAndLock(fm.DaoObjectFile, objectFileId)

		if err != nil {

			return nil, err
		}
		for {

			object, err := dao.objectIO.ReadSizePrefixed(table)

			if err != nil {

				if errors.Is(err, io.EOF) {

					break
				}
				return nil, err
			}
			objects = append(objects, object)
		}
		dao.fileManager.CloseAndUnlock(table, &err)
	}
	return objects, err
}
