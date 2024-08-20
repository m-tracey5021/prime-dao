package queue

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}

type IProcessor[T any] interface {
	Process() IResult[T]
}
