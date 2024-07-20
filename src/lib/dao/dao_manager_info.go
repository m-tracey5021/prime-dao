package dao

type DaoManagerInfo struct {
	fileNames map[uint64]string

	availableForWriting uint64
}
