package dao

import (
	"errors"
	"io"
	"os"
	"sync"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/ht"
	"transformer/src/lib/data/schema"
)

type IDaoMetadataManager[T schema.Identifiable] interface {
	AssignId(object T) T

	AvailableTable() uint64

	// AddAvailableTable(tableId uint64)

	// RemoveAvailableTable(tableId uint64)

	GetIndex(id uint64) (*DaoIndex, error)

	GetAllIndexes() ([]*DaoIndex, error)

	GetMetadata(id uint64) (*DaoObjFile, error)

	DeleteIndex(id uint64) error

	DeleteMetadata(id uint64) error

	// UpdateAvailableTableForSave() error

	UpdateForSave(tableId uint64, object T, position uint64) error

	UpdateForDeletion(table *os.File, index DaoIndex) error

	UpdateIndexes(table *os.File, fileId, filePosition uint64) error
}

type DaoMetadataManager[T schema.Identifiable] struct {
	daoId uint64

	fileManager fm.IFileManager

	managingInfo *DaoMetadata[T]

	managingInfoIO dataio.IDataIO[DaoMetadata[T]]

	indexHashTable ht.ITSFHashTable[DaoIndex]

	objFileHashTable ht.ITSFHashTable[DaoObjFile]

	objectIO dataio.IDataIO[T]

	mu sync.Mutex
}

func NewMetadataManager[T schema.Identifiable](daoId uint64, fileManager fm.IFileManager) (*DaoMetadataManager[T], error) {

	managingInfoIO := dataio.DataIO[DaoMetadata[T]]{}

	indexHashTable := ht.FromFileManager[DaoIndex](daoId, fileManager.Concatenate("idx"))

	objFileHashTable := ht.FromFileManager[DaoObjFile](daoId, fileManager.Concatenate("obj_f"))

	managingFile, err := fileManager.OpenAndLock(fm.DaoManagingFile, daoId)

	defer fileManager.CloseAndUnlock(managingFile, &err)

	if err != nil {

		return &DaoMetadataManager[T]{}, err
	}
	size, err := fileManager.Size(managingFile)

	if err != nil {

		return &DaoMetadataManager[T]{}, err
	}
	var managingInfo DaoMetadata[T]

	if size > 0 {

		if managingInfo, err = managingInfoIO.ReadSizePrefixed(managingFile); err != nil {

			return &DaoMetadataManager[T]{}, err
		}

	} else {

		initialTable := uint64(0)

		initialObjFile := DaoObjFile{initialTable, 0}

		managingInfo = DaoMetadata[T]{

			ObjectIdCache: DaoIdCache{},

			ObjectCache: DaoCache[T]{},

			TableIdCache: DaoIdCache{&initialTable, make([]uint64, 0)},

			MaxObjects: uint64(2),

			AvailableTable: initialTable,

			FirstAvailableTable: []uint64{initialTable},
		}
		if _, err := managingInfoIO.WriteSizePrefixed(managingFile, managingInfo); err != nil {

			return &DaoMetadataManager[T]{}, err
		}
		objFileHashTable.Save(initialObjFile)
	}
	return &DaoMetadataManager[T]{

			daoId: daoId,

			fileManager: fileManager,

			managingInfo: &managingInfo,

			managingInfoIO: managingInfoIO,

			indexHashTable: indexHashTable,

			objFileHashTable: objFileHashTable,

			objectIO: dataio.DataIO[T]{},
		},
		err
}

func (manager *DaoMetadataManager[T]) GetIndex(id uint64) (*DaoIndex, error) {

	return manager.indexHashTable.Get(id)
}

func (manager *DaoMetadataManager[T]) GetAllIndexes() ([]*DaoIndex, error) {

	ids := make([]uint64, 0)

	for id := range *manager.managingInfo.ObjectIdCache.Cached + 1 {

		ids = append(ids, id)
	}
	return manager.indexHashTable.GetSome(ids...)
}

func (manager *DaoMetadataManager[T]) GetMetadata(id uint64) (*DaoObjFile, error) {

	return manager.objFileHashTable.Get(id)
}

func (manager *DaoMetadataManager[T]) DeleteIndex(id uint64) error {

	return manager.indexHashTable.Delete(id)
}

func (manager *DaoMetadataManager[T]) DeleteMetadata(id uint64) error {

	return manager.objFileHashTable.Delete(id)
}

func (manager *DaoMetadataManager[T]) UpdateIndexes(table *os.File, fileId, filePosition uint64) error {

	if err := manager.fileManager.GoTo(int(filePosition), table); err != nil {

		return err
	}
	for {

		currentPosition, err := manager.fileManager.CurrentPosition(table)

		if err != nil {

			return err
		}
		readObject, err := manager.objectIO.ReadSizePrefixed(table)

		if err != nil {

			if errors.Is(err, io.EOF) {

				break
			}
			return err
		}
		updatedIndex := DaoIndex{readObject.Id(), fileId, uint64(currentPosition)}

		if err := manager.indexHashTable.Update(updatedIndex); err != nil {

			return err
		}
	}
	return nil
}

// func (manager *DaoMetadataManager[T]) UpdateAvailableTableForSave() error {

// 	metadata, err := manager.objFileHashTable.Get(manager.AvailableTable())

// 	if err != nil {

// 		return err
// 	}
// 	if metadata.objectsWritten == manager.managingInfo.MaxObjects {

// 		tableId := manager.NewTableId()

// 		newMetadata := DaoObjFile{tableId, 0}

// 		if err := manager.objFileHashTable.Save(newMetadata); err != nil {

// 			return err
// 		}
// 		manager.managingInfo.AvailableTable = tableId
// 	}
// 	return err
// }

func (manager *DaoMetadataManager[T]) UpdateMetadataForSave(tableIdSavedTo uint64) error {

	metadata, err := manager.objFileHashTable.Get(tableIdSavedTo)

	if err != nil {

		return err
	}
	objectsWrittenAfterSave := metadata.objectsWritten + 1

	if objectsWrittenAfterSave == manager.managingInfo.MaxObjects {

		tableId := manager.NewTableId()

		newMetadata := DaoObjFile{tableId, 0}

		if err := manager.objFileHashTable.Save(newMetadata); err != nil {

			return err
		}
		manager.RemoveAvailableTable(tableIdSavedTo)

		manager.AddAvailableTable(tableId)
	}
	metadata.objectsWritten += 1

	if err := manager.objFileHashTable.Update(*metadata); err != nil {

		return err
	}
	return err
}

func (manager *DaoMetadataManager[T]) UpdateMetadataForDeletion(fileIdDeletedFrom uint64) (bool, error) {

	removeTable := false

	metadata, err := manager.objFileHashTable.Get(fileIdDeletedFrom)

	if err != nil {

		return removeTable, err
	}
	if metadata.objectsWritten == 1 {

		removeTable = true

		manager.DeleteMetadata(metadata.id)

		manager.DeleteTableId(fileIdDeletedFrom)

		manager.RemoveAvailableTable(fileIdDeletedFrom)

	} else {

		metadata.objectsWritten -= 1

		manager.objFileHashTable.Update(*metadata)

		manager.AddAvailableTable(fileIdDeletedFrom)
	}
	return removeTable, err
}

func (manager *DaoMetadataManager[T]) UpdateForSave(tableIdSavedTo uint64, object T, position uint64) error {

	index := DaoIndex{object.Id(), tableIdSavedTo, position}

	if err := manager.indexHashTable.Save(index); err != nil {

		return err
	}
	if err := manager.UpdateMetadataForSave(tableIdSavedTo); err != nil {

		return err
	}
	if err := manager.SaveManagingInfo(); err != nil {

		return err
	}
	return nil
}

func (manager *DaoMetadataManager[T]) UpdateForDeletion(table *os.File, index DaoIndex) error {

	manager.managingInfo.ObjectIdCache.DeleteId(index.id, &manager.mu)

	shouldDelete, err := manager.UpdateMetadataForDeletion(index.fileId)

	if err != nil {

		return err
	}
	if shouldDelete {

		if err := manager.fileManager.Remove(table); err != nil {

			return err
		}
	}
	if err := manager.UpdateIndexes(table, index.fileId, index.filePosition); err != nil {

		return err
	}
	// could run these two in parallel
	if err := manager.SaveManagingInfo(); err != nil {

		return err
	}
	return nil
}

func (manager *DaoMetadataManager[T]) SaveManagingInfo() error {

	managingFile, err := manager.fileManager.OpenAndLock(fm.DaoManagingFile, manager.daoId)

	defer manager.fileManager.CloseAndUnlock(managingFile, &err)

	_, err = manager.managingInfoIO.WriteSizePrefixed(managingFile, *manager.managingInfo)

	if err != nil {

		return err
	}
	return err
}
