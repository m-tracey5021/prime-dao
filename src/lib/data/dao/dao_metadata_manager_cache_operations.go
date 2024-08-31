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

	if len(manager.managingInfo.FirstAvailableTable) == 0 {

		return uint64(0)
	}
	return manager.managingInfo.FirstAvailableTable[0]
}

func (manager *DaoMetadataManager[T]) AddAvailableTable(tableId uint64) {

	// this should be a set so that you cant add the same table twice
	manager.managingInfo.FirstAvailableTable = append(manager.managingInfo.FirstAvailableTable, tableId)
}

func (manager *DaoMetadataManager[T]) RemoveAvailableTable(tableId uint64) {

	manager.managingInfo.FirstAvailableTable = slices.DeleteFunc(manager.managingInfo.FirstAvailableTable, func(element uint64) bool {

		return element == tableId
	})
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
