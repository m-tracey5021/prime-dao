package dao

import (
	"sync"
)

type DaoIdStore struct {
	Last *uint64

	Deleted []uint64
}

func (store *DaoIdStore) NewId(mu *sync.Mutex) uint64 {

	mu.Lock()

	if len(store.Deleted) > 0 {

		popped := store.Deleted[len(store.Deleted)-1] // Get the last element

		store.Deleted = store.Deleted[:len(store.Deleted)-1]

		return popped
	}
	if store.Last != nil {

		*store.Last += 1

	} else {

		initial := uint64(0)

		store.Last = &initial
	}
	mu.Unlock()

	return *store.Last
}

func (store *DaoIdStore) Current() uint64 {

	return *store.Last
}

func (store *DaoIdStore) DeleteId(id uint64, mu *sync.Mutex) {

	mu.Lock()

	store.Deleted = append(store.Deleted, id)

	mu.Unlock()
}
