package dao

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

	return manager.managingInfo.ObjectCache.NewId(&manager.mu)
}

func (manager *DaoMetadataManager[T]) NewTableId() uint64 {

	return manager.managingInfo.TableCache.NewId(&manager.mu)
}

func (manager *DaoMetadataManager[T]) DeleteObjectId(id uint64) {

	manager.managingInfo.ObjectCache.DeleteId(id, &manager.mu)
}

func (manager *DaoMetadataManager[T]) DeleteTableId(id uint64) {

	manager.managingInfo.TableCache.DeleteId(id, &manager.mu)
}

func (manager *DaoMetadataManager[T]) AvailableTable() uint64 {

	return manager.managingInfo.AvailableTable
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
