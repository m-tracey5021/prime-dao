package dao

import (
	"slices"
	"sync"
)

type DaoIdStore struct {
	Last int

	Deleted []uint64
}

func NewIdStore() DaoIdStore {

	return DaoIdStore{-1, []uint64{}}
}

func (store *DaoIdStore) NewId(mu *sync.Mutex) uint64 {

	mu.Lock()

	defer mu.Unlock()

	if len(store.Deleted) > 0 {

		popped := store.Deleted[len(store.Deleted)-1] // Get the last element

		store.Deleted = store.Deleted[:len(store.Deleted)-1]

		return popped
	}
	if store.Last == -1 {

		store.Last = 0

	} else {

		store.Last += 1
	}
	return uint64(store.Last)
}

func (store *DaoIdStore) Current() uint64 {

	return uint64(store.Last)
}

func (store *DaoIdStore) AllIds() []uint64 {

	ids := []uint64{}

	for id := range uint64(store.Last) + 1 {

		if !slices.Contains(store.Deleted, id) {

			ids = append(ids, id)
		}
	}
	return ids
}

func (store *DaoIdStore) DeleteId(id uint64, mu *sync.Mutex) {

	mu.Lock()

	store.Deleted = append(store.Deleted, id)

	mu.Unlock()
}
