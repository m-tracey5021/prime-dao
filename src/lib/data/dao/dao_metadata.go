package dao

import "prime-dao/src/lib/data/schema"

type DaoMetadata[T schema.Identifiable] struct {
	ObjectIdStore DaoIdStore

	TableIdStore DaoIdStore

	Cache DaoCache[T]

	MaxObjects uint64

	AvailableTableObjectCount map[uint64]int
}
