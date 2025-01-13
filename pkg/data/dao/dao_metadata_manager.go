package dao

import (
	"errors"
	"io"
	"os"
	"prime-dao/pkg/data/dataio"
	"prime-dao/pkg/data/fm"
	"prime-dao/pkg/data/ht"
	"prime-dao/pkg/data/schema"
	"sync"
)

type IDaoMetadataManager[T schema.Identifiable] interface {
	NewId() uint64

	AllObjectIds() []uint64

	GetIndex(id uint64) (*DaoIndex, error)

	DeleteIndex(id uint64) error

	GetCached(id uint64) *T

	SaveToCache(object T)

	UpdateCache(object T)

	DeleteFromCache(id uint64)

	UpdateIndexes(table *os.File, fileId, filePosition uint64) error

	UpdateMetadataPreSave() uint64

	UpdateMetadataPostSave(tableId uint64, object T, position uint64) error

	UpdateMetadataForDeletion(table *os.File, index DaoIndex) error

	SaveManagingInfo() error
}

type DaoMetadataManager[T schema.Identifiable] struct {
	daoId uint64

	fileManager fm.IFileManager

	managingInfo *DaoMetadata[T]

	managingInfoIO dataio.IDataIO[DaoMetadata[T]]

	indexHashTable ht.ITSFHashTable[DaoIndex]

	objectIO dataio.IDataIO[T]

	idMu sync.Mutex

	cacheMu sync.Mutex

	mu sync.Mutex
}

func NewMetadataManager[T schema.Identifiable](daoId uint64, fileManager fm.IFileManager) (*DaoMetadataManager[T], error) {

	managingInfoIO := dataio.DataIO[DaoMetadata[T]]{}

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

		tableMapping := make(map[uint64]int)

		tableMapping[initialTable] = 0

		managingInfo = DaoMetadata[T]{

			ObjectIdStore: DaoIdStore{},

			TableIdStore: DaoIdStore{&initialTable, make([]uint64, 0)},

			Cache: NewCache[T](),

			MaxObjects: uint64(2),

			AvailableTableObjectCount: tableMapping,
		}
		if _, err := managingInfoIO.WriteSizePrefixed(managingFile, managingInfo); err != nil {

			return &DaoMetadataManager[T]{}, err
		}
	}
	return &DaoMetadataManager[T]{

			daoId: daoId,

			fileManager: fileManager,

			managingInfo: &managingInfo,

			managingInfoIO: managingInfoIO,

			indexHashTable: ht.FromFileManager[DaoIndex](daoId, fileManager.Concatenate("idx")),

			objectIO: dataio.DataIO[T]{},
		},
		err
}

func (manager *DaoMetadataManager[T]) GetIndex(id uint64) (*DaoIndex, error) {

	return manager.indexHashTable.Get(id)
}

func (manager *DaoMetadataManager[T]) DeleteIndex(id uint64) error {

	return manager.indexHashTable.Delete(id)
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

func (manager *DaoMetadataManager[T]) UpdateMetadataPreSave() uint64 {

	manager.mu.Lock()

	for table, objectCount := range manager.managingInfo.AvailableTableObjectCount {

		if objectCount < int(manager.managingInfo.MaxObjects) {

			manager.managingInfo.AvailableTableObjectCount[table] += 1

			manager.mu.Unlock()

			return table
		}
	}
	newTableId := manager.NewTableId()

	manager.managingInfo.AvailableTableObjectCount[newTableId] = 1

	manager.mu.Unlock()

	return newTableId
}

func (manager *DaoMetadataManager[T]) UpdateMetadataPostSave(tableIdSavedTo uint64, object T, position uint64) error {

	index := DaoIndex{object.Id(), tableIdSavedTo, position}

	if err := manager.indexHashTable.Save(index); err != nil {

		return err
	}
	return nil
}

func (manager *DaoMetadataManager[T]) UpdateMetadataForDeletion(table *os.File, index DaoIndex) error {

	manager.managingInfo.ObjectIdStore.DeleteId(index.id, &manager.mu)

	manager.mu.Lock()

	count, ok := manager.managingInfo.AvailableTableObjectCount[index.fileId]

	if ok {

		if count == 1 {

			if err := manager.fileManager.Remove(table); err != nil {

				return err
			}
			delete(manager.managingInfo.AvailableTableObjectCount, index.fileId)

		} else {

			manager.managingInfo.AvailableTableObjectCount[index.fileId] -= 1
		}
	}
	manager.mu.Unlock()

	if err := manager.UpdateIndexes(table, index.fileId, index.filePosition); err != nil {

		return err
	}
	return nil
}
