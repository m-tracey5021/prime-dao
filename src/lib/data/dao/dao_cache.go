package dao

import "transformer/src/lib/data/schema"

type DaoCache[T schema.Identifiable] struct {
	Cached map[uint64]T

	MaxSize int
}

func (cache DaoCache[T]) Get(id uint64) *T {

	object, ok := cache.Cached[id]

	if ok {

		return &object
	}
	return nil
}
