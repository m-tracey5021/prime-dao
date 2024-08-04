package dao

type DaoIdentifiable[T any] struct {
	Id uint64

	Object T
}
