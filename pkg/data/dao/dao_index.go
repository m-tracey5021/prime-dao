package dao

import (
	"encoding/binary"
	"os"
	"prime-dao/pkg/data/schema"
)

type DaoIndex struct {
	id uint64

	fileId uint64

	filePosition uint64
}

func (index DaoIndex) Id() uint64 {

	return index.id
}
func (index DaoIndex) SetId(id uint64) schema.Identifiable {

	return DaoIndex{id, index.fileId, index.filePosition}
}

func (index DaoIndex) Size() int {

	return 24
}

func (index DaoIndex) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, index.id); err != nil {

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

	var indexId uint64

	var fileId uint64

	var filePosition uint64

	if err := binary.Read(file, binary.LittleEndian, &indexId); err != nil {

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
