package dao

import (
	"slices"
	"transformer/src/lib/data/schema"
)

type DaoCache[T schema.Identifiable] struct {
	Cached map[uint64]T

	LeastRecentlyUsed []uint64

	MaxSize int
}

func (cache *DaoCache[T]) Get(id uint64) *T {

	object, ok := cache.Cached[id]

	if ok {

		cache.LeastRecentlyUsed = slices.DeleteFunc(cache.LeastRecentlyUsed, func(element uint64) bool {

			return element == id
		})
		cache.LeastRecentlyUsed = append(cache.LeastRecentlyUsed, id)

		return &object
	}
	return nil
}

func (cache *DaoCache[T]) Save(object T) {

	if len(cache.Cached) == cache.MaxSize {

		last := len(cache.LeastRecentlyUsed)

		delete(cache.Cached, cache.LeastRecentlyUsed[last])
	}
	cache.Cached[object.Id()] = object

	cache.LeastRecentlyUsed = append(cache.LeastRecentlyUsed, object.Id())
}
