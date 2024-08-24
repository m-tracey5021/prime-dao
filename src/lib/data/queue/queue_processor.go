package queue

import "sync"

type IProcessableRequest[T any] interface {
	RequestId() uint64

	ObjectId() uint64

	Dependencies() chan RequestContext

	Process(*sync.WaitGroup, *QueueContext[T]) IResult[T]
}

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}
