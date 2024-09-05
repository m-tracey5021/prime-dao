package dao

import (
	"tsf-dao/src/lib/data/queue"
	"tsf-dao/src/lib/data/schema"
)

type DaoTransaction[T schema.Identifiable] struct {
	requestFactory *DaoRequestFactory[T]

	requests []queue.IProcessableRequest[T]

	dependencyResolver *queue.QueueDependencyResolver[T]

	dependencies []uint64
}

func (transaction *DaoTransaction[T]) WithDependency(requestIds ...uint64) *DaoTransaction[T] {

	transaction.dependencies = requestIds

	return transaction
}

func (transaction *DaoTransaction[T]) Save(object T) uint64 {

	request := transaction.requestFactory.CreateSaveRequest(object, len(transaction.dependencies))

	transaction.requests = append(transaction.requests, request)

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Get(objectId uint64) uint64 {

	request := transaction.requestFactory.CreateGetRequest(objectId, len(transaction.dependencies))

	transaction.requests = append(transaction.requests, request)

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Update(object T) uint64 {

	request := transaction.requestFactory.CreateUpdateRequest(object, len(transaction.dependencies))

	transaction.requests = append(transaction.requests, request)

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Delete(objectId uint64) uint64 {

	request := transaction.requestFactory.CreateDeleteRequest(objectId, len(transaction.dependencies))

	transaction.requests = append(transaction.requests, request)

	transaction.dependencyResolver.AddDependency(request, transaction.dependencies)

	transaction.dependencies = []uint64{}

	return request.RequestId()
}
