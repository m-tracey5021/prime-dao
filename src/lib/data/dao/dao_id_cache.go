package dao

import (
	"sync"
)

type DaoIdCache struct {
	Cached *uint64

	Deleted []uint64
}

func (cache *DaoIdCache) NewId(mu *sync.Mutex) uint64 {

	mu.Lock()

	if len(cache.Deleted) > 0 {

		popped := cache.Deleted[len(cache.Deleted)-1] // Get the last element

		cache.Deleted = cache.Deleted[:len(cache.Deleted)-1]

		return popped
	}
	if cache.Cached != nil {

		*cache.Cached += 1

	} else {

		initial := uint64(0)

		cache.Cached = &initial
	}
	mu.Unlock()

	return *cache.Cached
}

func (cache *DaoIdCache) Current() uint64 {

	return *cache.Cached
}

func (cache *DaoIdCache) DeleteId(id uint64, mu *sync.Mutex) {

	mu.Lock()

	cache.Deleted = append(cache.Deleted, id)

	mu.Unlock()
}
