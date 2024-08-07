package ht

import (
	"encoding/binary"
	"os"
	"transformer/src/lib/dao/schema"
)

type HashTableCollisionTableHeader struct {
	// Used to check whether or not to initialise another table
	objectsWritten uint64
}

func (header HashTableCollisionTableHeader) Size() int {

	return 8
}

func (header HashTableCollisionTableHeader) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, header.objectsWritten); err != nil {

		return err
	}
	return nil
}

func (header HashTableCollisionTableHeader) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var objectsWritten uint64

	if err := binary.Read(file, binary.LittleEndian, &objectsWritten); err != nil {

		return nil, err
	}
	return HashTableCollisionTableHeader{objectsWritten}, nil
}
