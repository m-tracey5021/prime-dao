package dao

import (
	"errors"
	"io"
	"strconv"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
)

func (dao *Dao[T, U]) Get(objectId uuid.UUID) (*T, error) {

	cached := dao.metadata.Cache.Get(objectId, &dao.cacheMutex)

	if cached != nil {

		return cached, nil
	}
	index, err := dao.indexHashTable.Get(objectId)

	if err != nil {

		return nil, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoObjectFile, strconv.FormatUint(index.fileId, 10))

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

func (dao *Dao[T, U]) GetAll() ([]T, error) {

	objects := []T{}

	var err error

	for _, objectFileId := range dao.metadata.TableIdStore.AllIds() {

		table, err := dao.fileManager.OpenAndLock(fm.DaoObjectFile, strconv.FormatUint(objectFileId, 10))

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

// func (dao *Dao[T, U]) Scan(object T) (*T, error) {

// 	// add caching here aswell, as above
// 	index, err := dao.indexHashTable.Get(object.Id())

// 	if err != nil {

// 		return nil, err
// 	}

// 	x, y, z, a, b := dao.indexBTree.Search(*index)
// }
