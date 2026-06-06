package dao

import (
	"slices"

	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type QueueDependencyResolver[T schema.Orderable[U], U constraints.Ordered] struct {
	dependentProcesses map[IProcessableRequest[T, U]][]uint64 // e.g. process 0 depends on processes 1, 2, 3

	completed chan uint64
}

func NewResolver[T schema.Orderable[U], U constraints.Ordered]() *QueueDependencyResolver[T, U] {

	return &QueueDependencyResolver[T, U]{

		dependentProcesses: make(map[IProcessableRequest[T, U]][]uint64),

		completed: make(chan uint64, 100),
	}
}

func (resolver *QueueDependencyResolver[T, U]) ListenForCompletedDependencies() {

	go func() {

		for completedDependency := range resolver.completed {

			for process, dependencies := range resolver.dependentProcesses {

				found := slices.Contains(dependencies, completedDependency)

				if found {

					process.Dependencies() <- completedDependency
				}
			}
		}
	}()
}

func (resolver *QueueDependencyResolver[T, U]) AddDependency(request IProcessableRequest[T, U], dependencies ...uint64) {

	_, found := resolver.dependentProcesses[request]

	if found {

		resolver.dependentProcesses[request] = append(resolver.dependentProcesses[request], dependencies...)

	} else {

		resolver.dependentProcesses[request] = dependencies
	}
}

func (resolver *QueueDependencyResolver[T, U]) RemoveDependencies(requestId uint64) {

	for request := range resolver.dependentProcesses {

		if request.RequestId() == requestId {

			delete(resolver.dependentProcesses, request)

			return
		}
	}
}
