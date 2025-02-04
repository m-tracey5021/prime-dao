package queue

type IProcessableRequest[T any] interface {
	RequestId() uint64

	ObjectId() uint64

	Dependencies() chan uint64

	SetNumDeps(int)

	GetNumDeps() int

	Process() IResult[T]
}

type IResult[T any] interface {
	Object() *T

	Size() int

	Error() error
}
