package queue

import (
	"slices"
	"sync"
)

type RequestContext struct {
	RequestId uint64

	RequestWaitGroup *sync.WaitGroup
}

type QueueContext[T any] struct {
	dependentProcesses map[IProcessableRequest[T]][]uint64 // e.g. process 0 depends on processes 1, 2, 3

	completed chan RequestContext
}

func NewContext[T any]() *QueueContext[T] {

	return &QueueContext[T]{

		dependentProcesses: make(map[IProcessableRequest[T]][]uint64),

		completed: make(chan RequestContext, 100),
	}
}

func (context *QueueContext[T]) ResolveDependencies() {

	go func() {

		for req := range context.completed {

			for process, dependencies := range context.dependentProcesses {

				found := slices.Contains(dependencies, req.RequestId)

				if found {

					process.Dependencies() <- req
				}
			}
		}
	}()

}

func (context *QueueContext[T]) AddDependency(request IProcessableRequest[T], dependencies []uint64) {

	_, found := context.dependentProcesses[request]

	if found {

		context.dependentProcesses[request] = append(context.dependentProcesses[request], dependencies...)

	} else {

		context.dependentProcesses[request] = dependencies
	}
}

// func (context *QueueContext) Complete(id uint64) {

// 	dependencies := context.dependentProcesses[id]

// 	if len(dependencies) > 0 {

// 		// for {

// 		// 	var completed RequestContext

// 		// 	context.completed -> completed
// 		// }

// 		for completedRequest := range context.completed {

// 			found := slices.Contains(dependencies, completedRequest.RequestId)

// 			if found {

// 				completedRequest.RequestWaitGroup.Wait()

// 				fmt.Printf("successfully waited for dependency %v, while processing request %v", completedRequest.RequestId, id)
// 			}
// 		}
// 	}
// }
