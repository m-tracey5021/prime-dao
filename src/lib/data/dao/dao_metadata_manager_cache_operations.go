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

	return manager.managingInfo.ObjectIdCache.NewId(&manager.idMu)
}

func (manager *DaoMetadataManager[T]) NewTableId() uint64 {

	return manager.managingInfo.TableIdCache.NewId(&manager.idMu)
}

func (manager *DaoMetadataManager[T]) DeleteObjectId(id uint64) {

	manager.managingInfo.ObjectIdCache.DeleteId(id, &manager.idMu)
}

func (manager *DaoMetadataManager[T]) DeleteTableId(id uint64) {

	manager.managingInfo.TableIdCache.DeleteId(id, &manager.idMu)
}

func (manager *DaoMetadataManager[T]) GetCached(id uint64) *T {

	return manager.managingInfo.ObjectCache.Get(id)
}

func (manager *DaoMetadataManager[T]) SaveToCache(object T) {

	manager.managingInfo.ObjectCache.Save(object)
}
