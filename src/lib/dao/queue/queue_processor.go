package queue

type RequestType int

const (
	SaveRequest = iota

	GetRequest

	UpdateRequest

	DeleteRequest
)

type IRequest[T any] interface {
	Type() RequestType

	RequestId() uint64

	ObjectId() uint64

	Object() *T
}

type IResult[T any] interface {
	Object() *T

	Error() error
}

type IQueueProcessor[T any] interface {
	Process(IRequest[T]) IResult[T]
}
