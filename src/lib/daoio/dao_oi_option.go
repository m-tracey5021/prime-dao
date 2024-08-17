package daoio

type DaoIOOption[T any] struct {
	value *T
}

func (option DaoIOOption[T]) HasValue() bool {

	return option.value != nil
}
