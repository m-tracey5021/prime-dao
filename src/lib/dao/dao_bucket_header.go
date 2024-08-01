package dao

import (
	"encoding/binary"
	"os"
	"transformer/src/lib/dao/schema"
)

type DaoBucketHeader struct {
	occupied bool

	// If the original item has been deleted, dont bother searching for a collision
	// deleted bool

	// This flag checks whether or not a collision table has been created prior
	previousCollision bool

	// This value corresponds to the fileId where the collisions are stored
	collisionTableId uint64
}

func NewHeader(occupied bool) DaoBucketHeader {

	return DaoBucketHeader{}
}

func (header DaoBucketHeader) Size() int {

	return 10
}

func (header DaoBucketHeader) WriteSelf(file *os.File) error {

	occupied := boolToByte(header.occupied)

	// deleted := boolToByte(header.deleted)

	previousCollision := boolToByte(header.previousCollision)

	if err := binary.Write(file, binary.LittleEndian, occupied); err != nil {

		return err
	}
	// if err := binary.Write(file, binary.LittleEndian, deleted); err != nil {

	// 	return err
	// }
	if err := binary.Write(file, binary.LittleEndian, previousCollision); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, header.collisionTableId); err != nil {

		return err
	}
	return nil
}

func (header DaoBucketHeader) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var occupied byte

	// var deleted byte

	var previousCollision byte

	var collisionTableId uint64

	if err := binary.Read(file, binary.LittleEndian, &occupied); err != nil {

		return nil, err
	}
	// if err := binary.Read(file, binary.LittleEndian, &deleted); err != nil {

	// 	return nil, err
	// }
	if err := binary.Read(file, binary.LittleEndian, &previousCollision); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &collisionTableId); err != nil {

		return nil, err
	}
	return DaoBucketHeader{byteToBool(occupied), byteToBool(previousCollision), collisionTableId}, nil
}
