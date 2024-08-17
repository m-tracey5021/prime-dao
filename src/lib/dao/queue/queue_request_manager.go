package queue

type QueueRequestManager[T any] struct {
	requestFactory IQueueRequestFactory[T]

	queue *TSFQueue[T]
}

func NewQueueRequestManager[T any](requestFactory IQueueRequestFactory[T], queue *TSFQueue[T]) QueueRequestManager[T] {

	return QueueRequestManager[T]{requestFactory, queue}
}

// =============== Save ===============

func (manager *QueueRequestManager[T]) Save(object T) error {

	requestId := manager.queue.NewRequestId()

	request := manager.requestFactory.CreateSaveRequest(requestId, object)

	return manager.queue.ProcessSync(request).Error()
}

func (manager *QueueRequestManager[T]) SaveAsync(objects ...T) {

	requests := make([]IRequest[T], 0)

	for _, object := range objects {

		requestId := manager.queue.NewRequestId()

		request := manager.requestFactory.CreateSaveRequest(requestId, object)

		requests = append(requests, request)
	}
	manager.queue.ProcessAsync(requests...)
}

// =============== Get ===============

func (manager *QueueRequestManager[T]) Get(id uint64) (*T, error) {

	requestId := manager.queue.NewRequestId()

	request := manager.requestFactory.CreateGetRequest(requestId, id)

	result := manager.queue.ProcessSync(request)

	return result.Object(), result.Error()
}

func (manager *QueueRequestManager[T]) GetAsync(ids ...uint64) {

	requests := make([]IRequest[T], 0)

	for _, id := range ids {

		requestId := manager.queue.NewRequestId()

		request := manager.requestFactory.CreateGetRequest(requestId, id)

		requests = append(requests, request)
	}
	manager.queue.ProcessAsync(requests...)
}

// =============== Update ===============

func (manager *QueueRequestManager[T]) Update(object T) error {

	requestId := manager.queue.NewRequestId()

	request := manager.requestFactory.CreateUpdateRequest(requestId, object)

	return manager.queue.ProcessSync(request).Error()
}

func (manager *QueueRequestManager[T]) UpdateAsync(objects ...T) {

	requests := make([]IRequest[T], 0)

	for _, object := range objects {

		requestId := manager.queue.NewRequestId()

		request := manager.requestFactory.CreateUpdateRequest(requestId, object)

		requests = append(requests, request)
	}
	manager.queue.ProcessAsync(requests...)
}

// =============== Delete ===============

func (manager *QueueRequestManager[T]) Delete(id uint64) error {

	requestId := manager.queue.NewRequestId()

	request := manager.requestFactory.CreateDeleteRequest(requestId, id)

	result := manager.queue.ProcessSync(request)

	return result.Error()
}

func (manager *QueueRequestManager[T]) DeleteAsync(ids ...uint64) {

	requests := make([]IRequest[T], 0)

	for _, id := range ids {

		requestId := manager.queue.NewRequestId()

		request := manager.requestFactory.CreateDeleteRequest(requestId, id)

		requests = append(requests, request)
	}
	manager.queue.ProcessAsync(requests...)
}
