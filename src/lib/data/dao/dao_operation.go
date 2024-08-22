package dao

import "transformer/src/lib/data/queue"

func (dao TSFDao[T]) ProcessSave(object T) (*T, int, error) {

	getRequest := dao.requestFactory.CreateSaveRequest(dao.id, dao.fileManager, dao.metadataManager, object)

	result := dao.queue.ProcessSync(getRequest)

	return result.Object(), result.Size(), result.Error()
}

func (dao TSFDao[T]) ProcessSaveAsync(objects ...T) {

	requests := make([]queue.IProcessableRequest[T], 0)

	for _, object := range objects {

		getRequest := dao.requestFactory.CreateSaveRequest(dao.id, dao.fileManager, dao.metadataManager, object)

		requests = append(requests, getRequest)
	}
	dao.queue.ProcessAsync(requests...)
}

func (dao TSFDao[T]) ProcessGet(id uint64) (*T, error) {

	getRequest := dao.requestFactory.CreateGetRequest(dao.id, dao.fileManager, dao.metadataManager, id)

	result := dao.queue.ProcessSync(getRequest)

	return result.Object(), result.Error()
}

func (dao TSFDao[T]) ProcessGetAsync(ids ...uint64) {

	requests := make([]queue.IProcessableRequest[T], 0)

	for _, id := range ids {

		getRequest := dao.requestFactory.CreateGetRequest(dao.id, dao.fileManager, dao.metadataManager, id)

		requests = append(requests, getRequest)
	}
	dao.queue.ProcessAsync(requests...)
}

func (dao TSFDao[T]) ProcessUpdate(object T) (int, error) {

	getRequest := dao.requestFactory.CreateUpdateRequest(dao.id, dao.fileManager, dao.metadataManager, object)

	result := dao.queue.ProcessSync(getRequest)

	return result.Size(), result.Error()
}

func (dao TSFDao[T]) ProcessUpdateAsync(objects ...T) {

	requests := make([]queue.IProcessableRequest[T], 0)

	for _, object := range objects {

		getRequest := dao.requestFactory.CreateUpdateRequest(dao.id, dao.fileManager, dao.metadataManager, object)

		requests = append(requests, getRequest)
	}
	dao.queue.ProcessAsync(requests...)
}

func (dao TSFDao[T]) ProcessDelete(id uint64) (int, error) {

	getRequest := dao.requestFactory.CreateDeleteRequest(dao.id, dao.fileManager, dao.metadataManager, id)

	result := dao.queue.ProcessSync(getRequest)

	return result.Size(), result.Error()
}

func (dao TSFDao[T]) ProcessDeleteAsync(ids ...uint64) {

	requests := make([]queue.IProcessableRequest[T], 0)

	for _, id := range ids {

		getRequest := dao.requestFactory.CreateDeleteRequest(dao.id, dao.fileManager, dao.metadataManager, id)

		requests = append(requests, getRequest)
	}
	dao.queue.ProcessAsync(requests...)
}
