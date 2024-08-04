package dao

type DaoManagingInfo struct {
	ObjectIds []uint64

	TableIds []uint64

	MaxObjects uint64

	AvailableTable uint64
}
