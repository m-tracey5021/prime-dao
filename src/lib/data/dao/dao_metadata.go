package dao

import "transformer/src/lib/data/schema"

type DaoMetadata[T schema.Identifiable] struct {
	ObjectIdCache DaoIdCache

	ObjectCache DaoCache[T]

	TableIdCache DaoIdCache

	MaxObjects uint64

	AvailableTable uint64

	FirstAvailableTable []uint64
}
