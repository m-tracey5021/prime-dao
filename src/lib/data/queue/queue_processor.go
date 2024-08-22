package queue

type IProcessableRequest[T any] interface {
	ObjectId() uint64

	Process() IResult[T]
}

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}
