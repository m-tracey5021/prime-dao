package dao

import "transformer/src/lib/data/schema"

type DaoMetadata[T schema.Identifiable] struct {
	ObjectIdCache DaoIdCache

	ObjectCache DaoCache[T]

	TableIdCache DaoIdCache

	MaxObjects uint64

	AvailableTables []uint64

	AvailableTableObjectCount map[uint64]int
}
