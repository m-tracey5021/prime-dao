package dao

import "slices"

type CacheAlteration int

const (
	AddObject = iota

	AddTable

	RemoveObject

	RemoveTable
)

func (manager *DaoMetadataManager[T]) AssignId(object T) T {

	return object.SetId(manager.NewObjectId()).(T)
}

func (manager *DaoMetadataManager[T]) NewObjectId() uint64 {

	return manager.managingInfo.ObjectIdCache.NewId(&manager.mu)
}

func (manager *DaoMetadataManager[T]) NewTableId() uint64 {

	return manager.managingInfo.TableIdCache.NewId(&manager.mu)
}

func (manager *DaoMetadataManager[T]) DeleteObjectId(id uint64) {

	manager.managingInfo.ObjectIdCache.DeleteId(id, &manager.mu)
}

func (manager *DaoMetadataManager[T]) DeleteTableId(id uint64) {

	manager.managingInfo.TableIdCache.DeleteId(id, &manager.mu)
}

func (manager *DaoMetadataManager[T]) AvailableTable() uint64 {

	manager.mu.Lock()

	var firstAvailable uint64

	if len(manager.managingInfo.AvailableTables) == 0 {

		firstAvailable = uint64(0)

		manager.managingInfo.AvailableTableObjectCount[firstAvailable] = 0

	} else {

		firstAvailable = manager.managingInfo.AvailableTables[0]
	}
	count, ok := manager.managingInfo.AvailableTableObjectCount[firstAvailable]

	if ok {

		if count+1 == int(manager.managingInfo.MaxObjects) {

			delete(manager.managingInfo.AvailableTableObjectCount, firstAvailable)

		} else {

			manager.managingInfo.AvailableTableObjectCount[firstAvailable] += 1
		}
	}
	manager.mu.Unlock()

	return firstAvailable
}

func (manager *DaoMetadataManager[T]) AddAvailableTable(tableId uint64) {

	// this should be a set so that you cant add the same table twice
	manager.mu.Lock()

	manager.managingInfo.AvailableTables = append(manager.managingInfo.AvailableTables, tableId)

	manager.mu.Unlock()
}

func (manager *DaoMetadataManager[T]) RemoveAvailableTable(tableId uint64) {

	manager.mu.Lock()

	manager.managingInfo.AvailableTables = slices.DeleteFunc(manager.managingInfo.AvailableTables, func(element uint64) bool {

		return element == tableId
	})
	manager.mu.Unlock()
}

// func (manager *DaoMetadataManager[T]) AlterCache(id uint64, alteration CacheAlteration) {

// 	manager.mu.Lock()

// 	switch alteration {

// 	case AddTable:

// 		manager.managingInfo.TableIds = append(manager.managingInfo.TableIds, id)

// 	case AddObject:

// 		manager.managingInfo.ObjectIds = append(manager.managingInfo.ObjectIds, id)

// 	case RemoveTable:

// 		manager.managingInfo.TableIds = data.RemoveId(id, manager.managingInfo.TableIds)

// 	case RemoveObject:

// 		manager.managingInfo.ObjectIds = data.RemoveId(id, manager.managingInfo.ObjectIds)
// 	}
// 	manager.mu.Unlock()
// }
