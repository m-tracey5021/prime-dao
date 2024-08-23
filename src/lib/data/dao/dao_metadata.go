package dao

type DaoMetadata struct {
	ObjectCache DaoIdCache

	TableCache DaoIdCache

	MaxObjects uint64

	AvailableTable uint64
}
