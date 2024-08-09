package ht

import (
	"encoding/binary"
	"os"
	"transformer/src/lib"
	"transformer/src/lib/dao/schema"
)

type HashTableBucketHeader struct {
	occupied bool

	// This flag checks whether or not this bucket has been deleted
	deleted bool

	// This value corresponds to the fileId where the collisions are stored
	collisionTableId uint64
}

func (header HashTableBucketHeader) Size() int {

	return 10
}

func (header HashTableBucketHeader) WriteSelf(file *os.File) error {

	occupied := lib.BoolToByte(header.occupied)

	deleted := lib.BoolToByte(header.deleted)

	if err := binary.Write(file, binary.LittleEndian, occupied); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, deleted); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, header.collisionTableId); err != nil {

		return err
	}
	return nil
}

func (header HashTableBucketHeader) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var occupied byte

	var deleted byte

	var collisionTableId uint64

	if err := binary.Read(file, binary.LittleEndian, &occupied); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &deleted); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &collisionTableId); err != nil {

		return nil, err
	}
	return HashTableBucketHeader{lib.ByteToBool(occupied), lib.ByteToBool(deleted), collisionTableId}, nil
}
