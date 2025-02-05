package dao

import (
	"slices"

	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type QueueDependencyResolver[T schema.DescribedIdentifiable] struct {
	dependentProcesses map[IProcessableRequest[T]][]uint64 // e.g. process 0 depends on processes 1, 2, 3

	completed chan uint64
}

func NewResolver[T schema.DescribedIdentifiable]() *QueueDependencyResolver[T] {

	return &QueueDependencyResolver[T]{

		dependentProcesses: make(map[IProcessableRequest[T]][]uint64),

		completed: make(chan uint64, 100),
	}
}

func (resolver *QueueDependencyResolver[T]) ListenForCompletedDependencies() {

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

func (resolver *QueueDependencyResolver[T]) AddDependency(request IProcessableRequest[T], dependencies ...uint64) {

	_, found := resolver.dependentProcesses[request]

	if found {

		resolver.dependentProcesses[request] = append(resolver.dependentProcesses[request], dependencies...)

	} else {

		resolver.dependentProcesses[request] = dependencies
	}
}

func (resolver *QueueDependencyResolver[T]) RemoveDependencies(requestId uint64) {

	for request := range resolver.dependentProcesses {

		if request.RequestId() == requestId {

			delete(resolver.dependentProcesses, request)

			return
		}
	}
}
