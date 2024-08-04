package dao

import (
	"encoding/binary"
	"os"
	"transformer/src/lib/dao/schema"
)

type DaoMetadata struct {
	id uint64

	objectsWritten uint64
}

func (metadata DaoMetadata) Id() uint64 {

	return metadata.id
}

func (metadata DaoMetadata) Size() int {

	return 16
}

func (metadata DaoMetadata) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, metadata.id); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, metadata.objectsWritten); err != nil {

		return err
	}
	return nil
}

func (metadata DaoMetadata) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var metadataId uint64

	var objectsWritten uint64

	if err := binary.Read(file, binary.LittleEndian, &metadataId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &objectsWritten); err != nil {

		return nil, err
	}
	return DaoMetadata{metadataId, objectsWritten}, nil
}
