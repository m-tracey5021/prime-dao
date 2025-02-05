package dao

import (
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/fm"
)

func (dao *Dao[T]) UpdateMetadataForDeletion(table *os.File, index DaoIndex) error {

	dao.metadata.ObjectIdStore.DeleteId(index.id, &dao.idMutex)

	dao.metadataMutex.Lock()

	count, ok := dao.metadata.AvailableTableObjectCount[index.fileId]

	if ok {

		if count == 1 {

			if err := dao.fileManager.Remove(table); err != nil {

				return err
			}
			delete(dao.metadata.AvailableTableObjectCount, index.fileId)

			dao.metadata.TableIdStore.DeleteId(index.fileId, &dao.idMutex)

		} else {

			dao.metadata.AvailableTableObjectCount[index.fileId] -= 1
		}
	}
	dao.metadataMutex.Unlock()

	if err := dao.UpdateIndexes(table, index.fileId, index.filePosition); err != nil {

		return err
	}
	return nil
}

func (dao *Dao[T]) Delete(objectId uint64) (int, error) {

	index, err := dao.indexHashTable.Get(objectId)

	if err != nil {

		return 0, err
	}
	table, err := dao.fileManager.OpenAndLock(fm.DaoObjectFile, index.fileId)

	defer dao.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return 0, err
	}
	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return 0, err
	}
	sizeDeleted, err := dao.objectIO.Delete(table)

	if err != nil {

		return 0, err
	}
	if err := dao.UpdateMetadataForDeletion(table, *index); err != nil {

		return 0, err
	}
	if err := dao.indexHashTable.Delete(index.id); err != nil {

		return 0, err
	}
	dao.metadata.Cache.Delete(objectId, &dao.cacheMutex)

	return sizeDeleted, err
}
