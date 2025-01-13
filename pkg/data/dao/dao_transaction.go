package dao

import (
	"prime-dao/pkg/data/queue"
	"prime-dao/pkg/data/schema"
)

type DaoTransaction[T schema.Identifiable] struct {
	requestFactory *DaoRequestFactory[T]

	requests map[uint64]queue.IProcessableRequest[T]

	dependencyResolver *queue.QueueDependencyResolver[T]

	dependencies []uint64
}

func (transaction *DaoTransaction[T]) AddDependency(requestId uint64, dependencies ...uint64) {

	request := transaction.requests[requestId]

	transaction.dependencyResolver.AddDependency(request, dependencies...)
}

func (transaction *DaoTransaction[T]) WithDependency(requestIds ...uint64) *DaoTransaction[T] {

	transaction.dependencies = requestIds

	return transaction
}

func (transaction *DaoTransaction[T]) Save(object T) uint64 {

	request := transaction.requestFactory.CreateSaveRequest(object, len(transaction.dependencies))

	transaction.requests[request.RequestId()] = request

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Get(objectId uint64) uint64 {

	request := transaction.requestFactory.CreateGetRequest(objectId, len(transaction.dependencies))

	transaction.requests[request.RequestId()] = request

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Update(object T) uint64 {

	request := transaction.requestFactory.CreateUpdateRequest(object, len(transaction.dependencies))

	transaction.requests[request.RequestId()] = request

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Delete(objectId uint64) uint64 {

	request := transaction.requestFactory.CreateDeleteRequest(objectId, len(transaction.dependencies))

	transaction.requests[request.RequestId()] = request

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) RemoveRequest(requestId uint64) {

	delete(transaction.requests, requestId)

	transaction.dependencyResolver.RemoveDependency(requestId)

	transaction.dependencies = []uint64{}
}
