package queue

type IQueueRequestFactory[T any] interface {
	CreateSaveRequest(uint64, T) IRequest[T]

	CreateGetRequest(uint64, uint64) IRequest[T]

	CreateUpdateRequest(uint64, T) IRequest[T]

	CreateDeleteRequest(uint64, uint64) IRequest[T]
}
