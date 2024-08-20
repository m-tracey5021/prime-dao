package req

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoRequestManager[T schema.Identifiable] struct {
	queue *queue.TSFQueue[T]
}

// func NewRequestManager[T schema.Identifiable](path, descriptor string, id uint64) (*queue.QueueRequestManager[T], error) {

// 	dao, err := dao.New[T](path, descriptor, id)

// 	if err != nil {

// 		return nil, err
// 	}
// 	processor := DaoRequestProcessor[T]{*dao}

// 	requestQueue := queue.NewQueue(10, 5, 100, processor)

// 	requestFactory := DaoRequestFactory[T]{}

// 	requestManager := queue.NewQueueRequestManager[T](requestFactory, requestQueue)

// 	return &requestManager, err
// }

// =============== Save ===============

func (manager *DaoRequestManager[T]) Save(object T) error {

	requestId := manager.queue.NewRequestId()

	request := DaoSaveRequest[T]{requestId, object}

	return manager.queue.ProcessSync(&request).Error()
}

func (manager *DaoRequestManager[T]) SaveAsync(objects ...T) {

	requests := make([]queue.IRequest[T], 0)

	for _, object := range objects {

		requestId := manager.queue.NewRequestId()

		requests = append(requests, &DaoSaveRequest[T]{requestId, object})
	}
	manager.queue.ProcessAsync(requests...)
}

// =============== Get ===============

func (manager *DaoRequestManager[T]) Get(id uint64) (*T, error) {

	requestId := manager.queue.NewRequestId()

	request := DaoGetRequest[T]{requestId, id}

	result := manager.queue.ProcessSync(&request)

	return result.Object(), result.Error()
}

func (manager *DaoRequestManager[T]) GetAsync(ids ...uint64) {

	requests := make([]queue.IRequest[T], 0)

	for _, id := range ids {

		requestId := manager.queue.NewRequestId()

		requests = append(requests, &DaoGetRequest[T]{requestId, id})
	}
	manager.queue.ProcessAsync(requests...)
}

// =============== Update ===============

func (manager *DaoRequestManager[T]) Update(object T) error {

	requestId := manager.queue.NewRequestId()

	request := DaoUpdateRequest[T]{requestId, object}

	return manager.queue.ProcessSync(&request).Error()
}

func (manager *DaoRequestManager[T]) UpdateAsync(objects ...T) {

	requests := make([]queue.IRequest[T], 0)

	for _, object := range objects {

		requestId := manager.queue.NewRequestId()

		requests = append(requests, &DaoUpdateRequest[T]{requestId, object})
	}
	manager.queue.ProcessAsync(requests...)
}

// =============== Delete ===============

func (manager *DaoRequestManager[T]) Delete(id uint64) error {

	requestId := manager.queue.NewRequestId()

	request := DaoDeleteRequest[T]{requestId, id}

	result := manager.queue.ProcessSync(&request)

	return result.Error()
}

func (manager *DaoRequestManager[T]) DeleteAsync(ids ...uint64) {

	requests := make([]queue.IRequest[T], 0)

	for _, id := range ids {

		requestId := manager.queue.NewRequestId()

		requests = append(requests, &DaoDeleteRequest[T]{requestId, id})
	}
	manager.queue.ProcessAsync(requests...)
}
