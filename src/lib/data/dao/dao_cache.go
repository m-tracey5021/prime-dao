package dao

import "transformer/src/lib/data/schema"

type DaoCache[T schema.Identifiable] struct {
	Cached map[uint64]T

	MaxSize int
}
