package ht

import (
	"encoding/binary"
	"os"
	"transformer/src/lib"
	"transformer/src/lib/dao/schema"
)

type HashTableBucketHeader struct {
	occupied bool

	// This flag checks whether or not a collision table has been created prior
	previousCollision bool

	// This value corresponds to the fileId where the collisions are stored
	collisionTableId uint64
}

func (header HashTableBucketHeader) Size() int {

	return 10
}

func (header HashTableBucketHeader) WriteSelf(file *os.File) error {

	occupied := lib.BoolToByte(header.occupied)

	previousCollision := lib.BoolToByte(header.previousCollision)

	if err := binary.Write(file, binary.LittleEndian, occupied); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, previousCollision); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, header.collisionTableId); err != nil {

		return err
	}
	return nil
}

func (header HashTableBucketHeader) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var occupied byte

	var previousCollision byte

	var collisionTableId uint64

	if err := binary.Read(file, binary.LittleEndian, &occupied); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &previousCollision); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &collisionTableId); err != nil {

		return nil, err
	}
	return HashTableBucketHeader{lib.ByteToBool(occupied), lib.ByteToBool(previousCollision), collisionTableId}, nil
}
