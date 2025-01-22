package dao

import "github.com/m-tracey5021/prime-dao/pkg/data/schema"

type DaoMetadata[T schema.Identifiable] struct {
	ObjectIdStore DaoIdStore

	TableIdStore DaoIdStore

	Cache DaoCache[T]

	MaxObjects uint64

	AvailableTableObjectCount map[uint64]int
}
