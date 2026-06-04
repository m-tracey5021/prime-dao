package bt

import (
	"encoding/binary"
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type BTreeNode struct {
	Keys uuid.UUID

	Children uuid.UUID // start location of node to read
}

func (node BTreeNode) Size() int {

	return 32
}

func (node BTreeNode) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.BigEndian, node.Keys); err != nil {

		return err
	}
	if err := binary.Write(file, binary.BigEndian, node.Children); err != nil {

		return err
	}
	return nil
}

func (node BTreeNode) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var keys uuid.UUID

	var children uuid.UUID

	if err := binary.Read(file, binary.BigEndian, &keys); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &children); err != nil {

		return nil, err
	}
	return BTreeNode{keys, children}, nil
}
