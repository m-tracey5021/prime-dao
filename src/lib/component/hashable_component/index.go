package hashableComponent

import "transformer/src/lib/dao/schema"

type Index struct {
	id uint64

	fileNumber uint64

	filePosition uint64
}

func NewIndex(id, fileNumber, filePosition uint64) Index {

	return Index{id, fileNumber, filePosition}
}

func (index Index) Id() uint64 {

	return index.id
}

func (index Index) FileDescriptor() string {

	return "IDX"
}
func (index Index) Size() int {

	return 24
}

func (index Index) Empty() schema.Hashable {

	return NewIndex(0, 0, 0)
}
