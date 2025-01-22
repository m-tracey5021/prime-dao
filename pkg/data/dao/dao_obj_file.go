package dao

import (
	"encoding/binary"
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type DaoObjFile struct {
	id uint64

	objectsWritten uint64
}

func (objFile DaoObjFile) Id() uint64 {

	return objFile.id
}

func (objFile DaoObjFile) SetId(id uint64) schema.Identifiable {

	return DaoObjFile{id, objFile.objectsWritten}
}

func (objFile DaoObjFile) Size() int {

	return 16
}

func (objFile DaoObjFile) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, objFile.id); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, objFile.objectsWritten); err != nil {

		return err
	}
	return nil
}

func (objFile DaoObjFile) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var objFileId uint64

	var objectsWritten uint64

	if err := binary.Read(file, binary.LittleEndian, &objFileId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &objectsWritten); err != nil {

		return nil, err
	}
	return DaoObjFile{objFileId, objectsWritten}, nil
}
