package dao

import (
	"slices"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type DaoTransaction[T schema.Orderable[U], U constraints.Ordered] struct {
	requestIds map[uint64]struct{}

	requests []IProcessableRequest[T, U]

	dependencyResolver *QueueDependencyResolver[T, U]

	dependencies []uint64
}

func (transaction *DaoTransaction[T, U]) NewRequestId() uint64 {

	var id uint64 = 0

	for {

		if _, idExists := transaction.requestIds[id]; idExists {

			id += 1

		} else {

			return id
		}
	}
}

func (transaction *DaoTransaction[T, U]) AddRequest(request IProcessableRequest[T, U]) uint64 {

	if len(transaction.dependencies) > 0 {

		request.SetNumDeps(len(transaction.dependencies))

		transaction.dependencyResolver.AddDependency(request, transaction.dependencies...)

		transaction.dependencies = []uint64{}
	}
	transaction.requests = append(transaction.requests, request)

	transaction.requestIds[request.RequestId()] = struct{}{}

	return request.RequestId()
}

func (transaction *DaoTransaction[T, U]) Save(object T) uint64 {

	return transaction.AddRequest(

		&DaoSaveRequest[T, U]{

			requestId: transaction.NewRequestId(),

			object: object,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T, U]) Get(objectId uuid.UUID) uint64 {

	return transaction.AddRequest(

		&DaoGetRequest[T, U]{

			requestId: transaction.NewRequestId(),

			objectId: objectId,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T, U]) Update(object T) uint64 {

	return transaction.AddRequest(

		&DaoUpdateRequest[T, U]{

			requestId: transaction.NewRequestId(),

			object: object,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T, U]) Delete(objectId uuid.UUID) uint64 {

	return transaction.AddRequest(

		&DaoDeleteRequest[T, U]{

			requestId: transaction.NewRequestId(),

			objectId: objectId,

			dependencies: make(chan uint64),
		},
	)
}

func (transaction *DaoTransaction[T, U]) WithDependency(requestIds ...uint64) *DaoTransaction[T, U] {

	transaction.dependencies = requestIds

	return transaction
}

func (transaction *DaoTransaction[T, U]) RemoveRequest(requestId uint64) {

	transaction.requests = slices.DeleteFunc(transaction.requests, func(element IProcessableRequest[T, U]) bool {

		return element.RequestId() == requestId
	})
	transaction.dependencyResolver.RemoveDependencies(requestId)

	transaction.dependencies = []uint64{}
}
