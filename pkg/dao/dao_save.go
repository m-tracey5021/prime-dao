package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/fm"
)

func (dao *Dao[T]) UpdateMetadataPreSave() uint64 {

	dao.metadataMutex.Lock()

	for table, objectCount := range dao.metadata.AvailableTableObjectCount {

		if objectCount < int(dao.metadata.MaxObjects) {

			dao.metadata.AvailableTableObjectCount[table] += 1

			dao.metadataMutex.Unlock()

			return table
		}
	}
	newTableId := dao.metadata.TableIdStore.NewId(&dao.idMutex)

	dao.metadata.AvailableTableObjectCount[newTableId] = 1

	dao.metadataMutex.Unlock()

	return newTableId
}

func (dao *Dao[T]) UpdateMetadataPostSave(tableIdSavedTo uint64, object T, position uint64) error {

	index := DaoIndex{object.Id(), tableIdSavedTo, position}

	if err := dao.indexHashTable.Save(index); err != nil {

		return err
	}
	return nil
}

func (dao *Dao[T]) Save(object T) (int, error) {

	availableTable := dao.UpdateMetadataPreSave()

	objectFile, err := dao.fileManager.OpenAndLock(fm.DaoObjectFile, availableTable)

	defer dao.fileManager.CloseAndUnlock(objectFile, &err)

	if err != nil {

		return 0, err
	}
	position, err := dao.fileManager.Size(objectFile)

	if err := dao.fileManager.GoTo(position, objectFile); err != nil {

		return 0, err
	}
	size, err := dao.objectIO.WriteSizePrefixed(objectFile, object)

	if err != nil {

		return 0, err
	}
	dao.UpdateMetadataPostSave(availableTable, object, uint64(position))

	return size, err
}
