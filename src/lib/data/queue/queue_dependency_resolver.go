package queue

import (
	"slices"
	"sync"
)

type QueueDependency struct {
	RequestId uint64

	RequestWaitGroup *sync.WaitGroup
}

type QueueDependencyResolver[T any] struct {
	dependentProcesses map[IProcessableRequest[T]][]uint64 // e.g. process 0 depends on processes 1, 2, 3

	completed chan QueueDependency
}

func NewResolver[T any]() *QueueDependencyResolver[T] {

	return &QueueDependencyResolver[T]{

		dependentProcesses: make(map[IProcessableRequest[T]][]uint64),

		completed: make(chan QueueDependency, 100),
	}
}

func (resolver *QueueDependencyResolver[T]) ListenForCompletedDependencies() {

	go func() {

		for completedDep := range resolver.completed {

			for process, dependencies := range resolver.dependentProcesses {

				found := slices.Contains(dependencies, completedDep.RequestId)

				if found {

					process.Dependencies() <- completedDep
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

func (resolver *QueueDependencyResolver[T]) RemoveDependency(requestId uint64) {

	for request := range resolver.dependentProcesses {

		if request.RequestId() == requestId {

			delete(resolver.dependentProcesses, request)

			return
		}
	}
}
