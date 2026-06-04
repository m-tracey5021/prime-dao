package ht

import (
	"encoding/binary"
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
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

	occupied := dataio.BoolToByte(header.occupied)

	deleted := dataio.BoolToByte(header.deleted)

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
	return HashTableBucketHeader{dataio.ByteToBool(occupied), dataio.ByteToBool(deleted), collisionTableId}, nil
}
