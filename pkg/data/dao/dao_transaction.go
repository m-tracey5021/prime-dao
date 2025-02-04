package dao

import (
	"slices"

	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
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

func (transaction *DaoTransaction[T]) AddRequest(request queue.IProcessableRequest[T]) uint64 {

	if len(transaction.dependencies) > 0 {

		request.SetNumDeps(len(transaction.dependencies))

		transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

		transaction.dependencies = []uint64{}
	}
	transaction.requests = append(transaction.requests, request)

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Save(object T) uint64 {

	request := transaction.requestFactory.CreateSaveRequest(object)

	return transaction.AddRequest(request)
}

func (transaction *DaoTransaction[T]) Get(objectId uint64) uint64 {

	request := transaction.requestFactory.CreateGetRequest(objectId)

	return transaction.AddRequest(request)
}

func (transaction *DaoTransaction[T]) Update(object T) uint64 {

	request := transaction.requestFactory.CreateUpdateRequest(object)

	return transaction.AddRequest(request)
}

func (transaction *DaoTransaction[T]) Delete(objectId uint64) uint64 {

	request := transaction.requestFactory.CreateDeleteRequest(objectId)

	return transaction.AddRequest(request)
}

func (transaction *DaoTransaction[T]) RemoveRequest(requestId uint64) {

	transaction.requests = slices.DeleteFunc(transaction.requests, func(element queue.IProcessableRequest[T]) bool {

		return element.RequestId() == requestId
	})
	transaction.dependencyResolver.RemoveDependencies(requestId)

	transaction.dependencies = []uint64{}
}
