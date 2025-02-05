package dao

import (
	"slices"

	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type DaoTransaction[T schema.Identifiable] struct {
	requestIds map[uint64]struct{}

	requests []IProcessableRequest[T]

	dependencyResolver *QueueDependencyResolver[T]

	dependencies []uint64
}

func (transaction *DaoTransaction[T]) NewRequestId() uint64 {

	var id uint64 = 0

	for {

		if _, idExists := transaction.requestIds[id]; idExists {

			id += 1

		} else {

			return id
		}
	}
}

func (transaction *DaoTransaction[T]) AddRequest(request IProcessableRequest[T]) uint64 {

	if len(transaction.dependencies) > 0 {

		request.SetNumDeps(len(transaction.dependencies))

		transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

		transaction.dependencies = []uint64{}
	}
	transaction.requests = append(transaction.requests, request)

	transaction.requestIds[request.RequestId()] = struct{}{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T]) Save(object T) uint64 {

	return transaction.AddRequest(

		&DaoSaveRequest[T]{

			requestId: transaction.NewRequestId(),

			object: object,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T]) Get(objectId uint64) uint64 {

	return transaction.AddRequest(

		&DaoGetRequest[T]{

			requestId: transaction.NewRequestId(),

			objectId: objectId,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T]) Update(object T) uint64 {

	return transaction.AddRequest(

		&DaoUpdateRequest[T]{

			requestId: transaction.NewRequestId(),

			object: object,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T]) Delete(objectId uint64) uint64 {

	return transaction.AddRequest(

		&DaoDeleteRequest[T]{

			requestId: transaction.NewRequestId(),

			objectId: objectId,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T]) WithDependency(requestIds ...uint64) *DaoTransaction[T] {

	transaction.dependencies = requestIds

	return transaction
}

func (transaction *DaoTransaction[T]) RemoveRequest(requestId uint64) {

	transaction.requests = slices.DeleteFunc(transaction.requests, func(element IProcessableRequest[T]) bool {

		return element.RequestId() == requestId
	})
	transaction.dependencyResolver.RemoveDependencies(requestId)

	transaction.dependencies = []uint64{}
}
