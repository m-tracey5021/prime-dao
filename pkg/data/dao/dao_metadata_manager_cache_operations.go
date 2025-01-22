package dao

import "github.com/m-tracey5021/prime-dao/pkg/data/fm"

func (manager *DaoMetadataManager[T]) NewId() uint64 {

	return manager.NewObjectId()
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

func (manager *DaoMetadataManager[T]) UpdateCache(object T) {

	manager.managingInfo.Cache.Update(object, &manager.cacheMu)
}

func (manager *DaoMetadataManager[T]) DeleteFromCache(id uint64) {

	manager.managingInfo.Cache.Delete(id, &manager.cacheMu)
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
