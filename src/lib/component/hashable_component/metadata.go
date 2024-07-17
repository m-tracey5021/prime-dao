package hashableComponent

type Metadata struct {
	id uint64

	objectsWritten uint64
}

func NewMetadata(id, objectsWritten uint64) Metadata {

	return Metadata{id, objectsWritten}
}

func (metadata Metadata) Id() uint64 {

	return metadata.id
}

func (metadata Metadata) FileDescriptor() string {

	return "MTD"
}
func (metadata Metadata) Size() int {

	return 16
}

func (metadata Metadata) Empty() Metadata {

	return NewMetadata(0, 0)
}
