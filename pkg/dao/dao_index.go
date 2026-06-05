package dao

import (
	"encoding/binary"
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type DaoIndex struct {
	id uuid.UUID

	fileId uint64

	filePosition uint64
}

func (index DaoIndex) Id() uuid.UUID {

	return index.id
}

func (index DaoIndex) Descriptor() string {

	return "idx"
}

func (index DaoIndex) Size() int {

	return 32
}

func (index DaoIndex) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.BigEndian, index.id); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, index.fileId); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, index.filePosition); err != nil {

		return err
	}
	return nil
}

func (index DaoIndex) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var indexId uuid.UUID

	var fileId uint64

	var filePosition uint64

	if err := binary.Read(file, binary.BigEndian, &indexId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &fileId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &filePosition); err != nil {

		return nil, err
	}
	return DaoIndex{indexId, fileId, filePosition}, nil
}

// func (index DaoIndex) Compare(other schema.Orderable) schema.Order {

// }
