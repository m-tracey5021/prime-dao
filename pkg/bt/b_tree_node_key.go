package bt

import (
	"encoding/binary"
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type BTreeNodeKey struct {
	Id uuid.UUID // 16 bytes

	FileId uint64 // 8 bytes

	FilePosition uint64 // 8 bytes
}

func (key BTreeNodeKey) Size() int {

	return 32
}

func (key BTreeNodeKey) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.BigEndian, key.Id); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, key.FileId); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, key.FilePosition); err != nil {

		return err
	}
	return nil
}

func (key BTreeNodeKey) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var id uuid.UUID

	var fileId uint64

	var filePosition uint64

	if err := binary.Read(file, binary.BigEndian, &id); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &fileId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &filePosition); err != nil {

		return nil, err
	}
	return BTreeNodeKey{id, fileId, filePosition}, nil
}
