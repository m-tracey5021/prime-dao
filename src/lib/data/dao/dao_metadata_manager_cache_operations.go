package dao

func (manager *DaoMetadataManager[T]) AssignId(object T) T {

	return object.SetId(manager.NewObjectId()).(T)
}

func (manager *DaoMetadataManager[T]) AllObjectIds() []uint64 {

	ids := make([]uint64, 0)

	for id := range *manager.managingInfo.ObjectIdStore.Last + 1 {

		ids = append(ids, id)
	}
	return ids
}

func (manager *DaoMetadataManager[T]) NewObjectId() uint64 {

	return manager.managingInfo.ObjectIdStore.NewId(&manager.idMu)
}

func (manager *DaoMetadataManager[T]) NewTableId() uint64 {

	return manager.managingInfo.TableIdStore.NewId(&manager.idMu)
}

func (manager *DaoMetadataManager[T]) DeleteObjectId(id uint64) {

	manager.managingInfo.ObjectIdStore.DeleteId(id, &manager.idMu)
}

func (manager *DaoMetadataManager[T]) DeleteTableId(id uint64) {

	manager.managingInfo.TableIdStore.DeleteId(id, &manager.idMu)
}

func (manager *DaoMetadataManager[T]) GetCached(id uint64) *T {

	return manager.managingInfo.Cache.Get(id, &manager.cacheMu)
}

func (manager *DaoMetadataManager[T]) SaveToCache(object T) {

	manager.managingInfo.Cache.Save(object, &manager.cacheMu)
}
