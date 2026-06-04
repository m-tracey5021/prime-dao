package bt

import (
	"bytes"
	"fmt"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/list"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type BTree[T schema.Identifiable] struct {
	order int

	fileManager fm.IFileManager

	root BTreeNode

	parentage map[BTreeNode]BTreeNode
}

func NewBTree[T schema.Identifiable](order int, fileManager fm.IFileManager) (BTree[T], error) {

	rootFile, err := fileManager.OpenAndLock(fm.BTree)

	defer fileManager.CloseAndUnlock(rootFile, &err)

	size, err := fileManager.Size(rootFile)

	if err != nil {

		return BTree[T]{}, err
	}
	if size == 0 {

		return BTree[T]{

				order: order,

				fileManager: fileManager,

				root: BTreeNode{uuid.Nil, uuid.Nil},

				parentage: map[BTreeNode]BTreeNode{},
			},
			err

	} else {

		nodeIO := dataio.FixedSizeDataIO[BTreeNode]{}

		root, err := nodeIO.Read(rootFile)

		if err != nil {

			return BTree[T]{}, err
		}
		return BTree[T]{order, fileManager, root, map[BTreeNode]BTreeNode{}}, err
	}
}

func (btree *BTree[T]) SaveRoot(node BTreeNode) error {

	btree.root = node

	rootFile, err := btree.fileManager.OpenAndLock(fm.BTree)

	defer btree.fileManager.CloseAndUnlock(rootFile, &err)

	nodeIO := dataio.FixedSizeDataIO[BTreeNode]{}

	err = nodeIO.Write(rootFile, node)

	return err
}

func (btree BTree[T]) IndexForKey(keys list.LinkedList[BTreeNodeKey], objectId uuid.UUID) (int, bool, error) {

	iterator, err := keys.Iter()

	if err != nil {

		return 0, false, err
	}
	defer iterator.Close(&err)

	key, err := iterator.Next()

	if err != nil {

		return 0, false, err
	}
	for key != nil {

		// logic

		uuidCompare := btree.CompareUUIDs(objectId, key.Id)

		if uuidCompare == -1 { // new < existing

			// new key points to the current compared key as the next key

			return iterator.Current(), false, err

		} else if uuidCompare == 1 { // new > existing

			// go to next key, unless this is last key

			key, err = iterator.Next()

			if err != nil {

				return 0, false, err
			}

		} else { // new == existing

			return iterator.Current(), true, err
		}
	}
	return int(keys.Size()), false, err
}

func (btree *BTree[T]) ToString() string {

	str := ""

	btree.BuildString(&str, btree.root, 0)

	return str
}

func (btree *BTree[T]) BuildString(str *string, node BTreeNode, level int) error {

	*str += "\n"

	for range level {

		*str += "\t"
	}
	*str += "["

	keys := list.From[BTreeNodeKey](btree.fileManager, node.Keys)

	keysActual, err := keys.ToSlice()

	if err != nil {

		return err
	}
	for _, key := range keysActual {

		*str += fmt.Sprintf("%v, ", key)
	}
	*str += "]"

	if node.Children != uuid.Nil {

		children := list.From[BTreeNode](btree.fileManager, node.Children)

		childrenActual, err := children.ToSlice()

		if err != nil {

			return err
		}
		for _, child := range childrenActual {

			btree.BuildString(str, child, level+1)
		}

	}
	return nil
}

func (btree BTree[T]) CompareUUIDs(a, b uuid.UUID) int {

	return bytes.Compare(a[:], b[:]) // Lexicographic comparison
}
