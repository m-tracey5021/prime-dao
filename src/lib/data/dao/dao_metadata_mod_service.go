package dao

import (
	"slices"
	"sync"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/schema"
)

type DaoMetadataModificationService[T schema.Identifiable] struct {
	daoId uint64

	fileManager fm.IFileManager

	managingInfo DaoMetadata[T]

	managingInfoIO dataio.IDataIO[DaoMetadata[T]]

	availableTables []uint64

	availableTableObjectCount map[uint64]int

	availableTableRequestCount int

	mu sync.Mutex
}

func (manager *DaoMetadataModificationService[T]) AssignId(object T) T {

	return object.SetId(manager.NewObjectId()).(T)
}

func (manager *DaoMetadataModificationService[T]) NewObjectId() uint64 {

	return manager.managingInfo.ObjectIdCache.NewId(&manager.mu)
}

func (manager *DaoMetadataModificationService[T]) NewTableId() uint64 {

	return manager.managingInfo.TableIdCache.NewId(&manager.mu)
}

func (manager *DaoMetadataModificationService[T]) DeleteObjectId(id uint64) {

	manager.managingInfo.ObjectIdCache.DeleteId(id, &manager.mu)
}

func (manager *DaoMetadataModificationService[T]) DeleteTableId(id uint64) {

	manager.managingInfo.TableIdCache.DeleteId(id, &manager.mu)
}

func (manager *DaoMetadataModificationService[T]) AvailableTablePreSave() uint64 {

	manager.mu.Lock()

	var firstAvailable uint64

	if len(manager.availableTables) == 0 {

		firstAvailable = uint64(0)

		manager.availableTableObjectCount[firstAvailable] = 0

	} else {

		firstAvailable = manager.availableTables[0]
	}
	count, ok := manager.availableTableObjectCount[firstAvailable]

	if ok {

		if count+1 == int(manager.managingInfo.MaxObjects) {

			delete(manager.availableTableObjectCount, firstAvailable)

		} else {

			manager.availableTableObjectCount[firstAvailable] += 1
		}
	}
	manager.mu.Unlock()

	return firstAvailable
}

func (manager *DaoMetadataModificationService[T]) AddAvailableTable(tableId uint64) {

	// this should be a set so that you cant add the same table twice
	manager.mu.Lock()

	manager.managingInfo.AvailableTables = append(manager.managingInfo.AvailableTables, tableId)

	manager.mu.Unlock()
}

func (manager *DaoMetadataModificationService[T]) RemoveAvailableTable(tableId uint64) {

	manager.mu.Lock()

	manager.managingInfo.AvailableTables = slices.DeleteFunc(manager.managingInfo.AvailableTables, func(element uint64) bool {

		return element == tableId
	})
	manager.mu.Unlock()
}

func (manager *DaoMetadataModificationService[T]) SaveManagingInfo() error {

	managingFile, err := manager.fileManager.OpenAndLock(fm.DaoManagingFile, manager.daoId)

	defer manager.fileManager.CloseAndUnlock(managingFile, &err)

	_, err = manager.managingInfoIO.WriteSizePrefixed(managingFile, manager.managingInfo)

	if err != nil {

		return err
	}
	return err
}
