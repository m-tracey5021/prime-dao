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

type DaoManagingInfo struct {
	ObjectIds []uint64

	TableIds []uint64

	MaxObjects uint64

	AvailableTable uint64
}

type IDaoMetadataManager[T schema.Identifiable] interface {
	AssignId(object T) T

	AvailableTable() uint64

	GetIndex(id uint64) (*DaoIndex, error)

	GetMetadata(id uint64) (*DaoObjFile, error)

	DeleteIndex(id uint64) error

	DeleteMetadata(id uint64) error

	UpdateForSave(table *os.File, object T, position uint64) error

	UpdateForDeletion(table *os.File, index DaoIndex) error

	UpdateIndexes(table *os.File, fileId, filePosition uint64) error
}

type DaoMetadataManager[T schema.Identifiable] struct {
	daoId uint64

	fileManager fm.IFileManager

	managingInfo DaoManagingInfo

	managingInfoIO dataio.IDataIO[DaoManagingInfo]

	indexHashTable ht.ITSFHashTable[DaoIndex]

	objFileHashTable ht.ITSFHashTable[DaoObjFile]

	objectIO dataio.IDataIO[T]

	mu sync.Mutex
}

func NewMetadataManager[T schema.Identifiable](daoId uint64, fileManager fm.IFileManager) (*DaoMetadataManager[T], error) {

	managingInfoIO := dataio.DataIO[DaoManagingInfo]{}

	managingFile, err := fileManager.OpenAndLock(fm.DaoManagingFile, daoId)

	defer fileManager.CloseAndUnlock(managingFile, &err)

	if err != nil {

		return &DaoMetadataManager[T]{}, err
	}
	size, err := fileManager.Size(managingFile)

	if err != nil {

		return &DaoMetadataManager[T]{}, err
	}
	var managingInfo *DaoManagingInfo

	if size > 0 {

		if managingInfo, err = managingInfoIO.ReadSizePrefixed(managingFile); err != nil {

			return &DaoMetadataManager[T]{}, err
		}

	} else {

		objectIds := make([]uint64, 0)

		fileIds := make([]uint64, 0)

		maxObjects := uint64(10) // Get from config

		available := uint64(0)

		managingInfo = &DaoManagingInfo{objectIds, fileIds, maxObjects, available}

		if _, err := managingInfoIO.WriteSizePrefixed(managingFile, *managingInfo); err != nil {

			return &DaoMetadataManager[T]{}, err
		}
	}
	return &DaoMetadataManager[T]{

			daoId: daoId,

			fileManager: fileManager,

			managingInfo: *managingInfo,

			managingInfoIO: dataio.DataIO[DaoManagingInfo]{},

			indexHashTable: ht.FromFileManager[DaoIndex](daoId, fileManager),

			objFileHashTable: ht.FromFileManager[DaoObjFile](daoId, fileManager),

			objectIO: dataio.DataIO[T]{},
		},
		err
}

func (manager *DaoMetadataManager[T]) GetIndex(id uint64) (*DaoIndex, error) {

	return manager.indexHashTable.Get(id)
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
		updatedIndex := DaoIndex{(*readObject).Id(), fileId, uint64(currentPosition)}

		if err := manager.indexHashTable.Update(updatedIndex); err != nil {

			return err
		}
	}
	return nil
}

func (manager *DaoMetadataManager[T]) UpdateMetadataForSave(fileIdSavedTo uint64) error {

	metadata, err := manager.objFileHashTable.Get(fileIdSavedTo)

	if err != nil {

		return err
	}
	objectsWrittenAfterSave := metadata.objectsWritten + 1

	if objectsWrittenAfterSave == manager.MaxObjects() {

		tableId := manager.NewTableId()

		manager.AlterCache(tableId, AddTable)

		manager.SetAvailableTable(tableId)
	}
	metadata.objectsWritten += 1

	if err := manager.objFileHashTable.Update(*metadata); err != nil {

		return err
	}
	return err
}

func (manager *DaoMetadataManager[T]) UpdateMetadataForDeletion(table *os.File, fileIdDeletedFrom uint64) error {

	metadata, err := manager.objFileHashTable.Get(fileIdDeletedFrom)

	if err != nil {

		return err
	}
	if metadata.objectsWritten == 1 {

		if err := manager.fileManager.Remove(table); err != nil {

			return err
		}
		manager.objFileHashTable.Delete(metadata.id)

		manager.AlterCache(fileIdDeletedFrom, RemoveTable)

	} else {

		metadata.objectsWritten -= 1

		manager.objFileHashTable.Update(*metadata)

		manager.SetAvailableTable(fileIdDeletedFrom)
	}
	return err
}

func (manager *DaoMetadataManager[T]) UpdateForSave(table *os.File, object T, position uint64) error {

	manager.AlterCache(object.Id(), AddObject)

	index := DaoIndex{object.Id(), manager.AvailableTable(), position}

	if err := manager.indexHashTable.Save(index); err != nil {

		return err
	}
	if err := manager.UpdateMetadataForSave(manager.AvailableTable()); err != nil {

		return err
	}
	if err := manager.SaveManagingInfo(); err != nil {

		return err
	}
	return nil
}

func (manager *DaoMetadataManager[T]) UpdateForDeletion(table *os.File, index DaoIndex) error {

	manager.AlterCache(index.id, RemoveObject)

	if err := manager.UpdateMetadataForDeletion(table, index.fileId); err != nil {

		return err
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

	_, err = manager.managingInfoIO.WriteSizePrefixed(managingFile, manager.managingInfo)

	if err != nil {

		return err
	}
	return err
}
