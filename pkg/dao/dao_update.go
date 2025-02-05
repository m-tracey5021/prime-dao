package dao

import (
	"errors"
	"io"
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/fm"
)

func (dao *Dao[T]) UpdateIndexes(objectFile *os.File, fileId, filePosition uint64) error {

	if err := dao.fileManager.GoTo(int(filePosition), objectFile); err != nil {

		return err
	}
	for {

		currentPosition, err := dao.fileManager.CurrentPosition(objectFile)

		if err != nil {

			return err
		}
		readObject, err := dao.objectIO.ReadSizePrefixed(objectFile)

		if err != nil {

			if errors.Is(err, io.EOF) {

				break
			}
			return err
		}
		updatedIndex := DaoIndex{readObject.Id(), fileId, uint64(currentPosition)}

		if err := dao.indexHashTable.Update(updatedIndex); err != nil {

			return err
		}
	}
	return nil
}

func (dao *Dao[T]) Update(object T) (int, error) {

	index, err := dao.indexHashTable.Get(object.Id())

	if err != nil {

		return 0, err
	}
	objectFile, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer dao.fileManager.CloseAndUnlock(objectFile, &err)

	if err != nil {

		return 0, err
	}
	if err := dao.fileManager.GoTo(int(index.filePosition), objectFile); err != nil {

		return 0, err
	}
	updatedSize, err := dao.objectIO.Update(objectFile, object)

	if err != nil {

		return 0, err
	}
	updatedFilePosition := int(index.filePosition) + updatedSize

	if err := dao.UpdateIndexes(objectFile, index.fileId, uint64(updatedFilePosition)); err != nil {

		return 0, err
	}
	dao.metadata.Cache.Update(object, &dao.cacheMutex)

	return updatedSize, err
}
