package queue

import "sync"

type IProcessableRequest[T any] interface {
	RequestId() uint64

	ObjectId() uint64

	Dependencies() chan QueueDependency

	SetNumDeps(int)

	ProcessWithDependencies(*sync.WaitGroup, *QueueDependencyResolver[T]) IResult[T]

	Process() IResult[T]
}

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}
