package dao

import (
	"slices"
	"sync"
	"tsf-dao/src/lib/data/schema"
)

type DaoCache[T schema.Identifiable] struct {
	Cached map[uint64]T

	LeastRecentlyUsed []uint64

	MaxSize int
}

func NewCache[T schema.Identifiable]() DaoCache[T] {

	return DaoCache[T]{
		Cached: make(map[uint64]T),

		LeastRecentlyUsed: make([]uint64, 0),

		MaxSize: 3,
	}
}

func (cache *DaoCache[T]) Get(id uint64, mu *sync.Mutex) *T {

	mu.Lock()

	object, ok := cache.Cached[id]

	if ok {

		cache.LeastRecentlyUsed = slices.DeleteFunc(cache.LeastRecentlyUsed, func(element uint64) bool {

			return element == id
		})
		cache.LeastRecentlyUsed = append(cache.LeastRecentlyUsed, id)

		mu.Unlock()

		return &object
	}
	mu.Unlock()

	return nil
}

func (cache *DaoCache[T]) Save(object T, mu *sync.Mutex) {

	mu.Lock()

	if len(cache.Cached) == cache.MaxSize {

		last := cache.LeastRecentlyUsed[0]

		delete(cache.Cached, last)

		cache.LeastRecentlyUsed = slices.Delete(cache.LeastRecentlyUsed, 0, 1)
	}
	cache.Cached[object.Id()] = object

	cache.LeastRecentlyUsed = append(cache.LeastRecentlyUsed, object.Id())

	mu.Unlock()
}
