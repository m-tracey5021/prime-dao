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

type BTree[T schema.FixedSizeIdentifiable, U schema.Orderable] struct {
	order int

	fileManager fm.IFileManager

	root BTreeNode

	parentage map[BTreeNode]BTreeNode

	// comparator func(a, b U) int

	getCompareValue func(keyId uuid.UUID) (*U, error)
}

func NewBTree[T schema.FixedSizeIdentifiable, U schema.Orderable](order int, fileManager fm.IFileManager, getCompareValue func(keyId uuid.UUID) (*U, error)) (BTree[T, U], error) {

	rootFile, err := fileManager.OpenAndLock(fm.BTree)

	defer fileManager.CloseAndUnlock(rootFile, &err)

	size, err := fileManager.Size(rootFile)

	if err != nil {

		return BTree[T, U]{}, err
	}
	if size == 0 {

		return BTree[T, U]{

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

			return BTree[T, U]{}, err
		}
		return BTree[T, U]{order, fileManager, root, map[BTreeNode]BTreeNode{}, getCompareValue}, err
	}
}

func (btree *BTree[T, U]) SaveRoot(node BTreeNode) error {

	btree.root = node

	rootFile, err := btree.fileManager.OpenAndLock(fm.BTree)

	defer btree.fileManager.CloseAndUnlock(rootFile, &err)

	nodeIO := dataio.FixedSizeDataIO[BTreeNode]{}

	err = nodeIO.Write(rootFile, node)

	return err
}

func (btree BTree[T, U]) IndexForExistingKey(keys list.FixedSizeLinkedList[T], targetKeyId uuid.UUID) (int, bool, error) {

	iterator, err := keys.Iter()

	if err != nil {

		return 0, false, err
	}
	defer iterator.Close(&err)

	currentKey, err := iterator.Next()

	if err != nil {

		return 0, false, err
	}
	for currentKey != nil {

		a, err := btree.getCompareValue(targetKeyId)

		b, err := btree.getCompareValue((*currentKey).Id())

		compareValue := (*a).Compare(*b)

		if compareValue == -1 { // new < existing

			// new key points to the current compared key as the next key

			return iterator.Current(), false, err

		} else if compareValue == 1 { // new > existing

			// go to next key, unless this is last key

			currentKey, err = iterator.Next()

			if err != nil {

				return 0, false, err
			}

		} else { // new == existing

			return iterator.Current(), true, err
		}
	}
	return int(keys.Size()), false, err
}

func (btree BTree[T, U]) IndexForNonExistingKey(keys list.FixedSizeLinkedList[T], object T, sortValue U) (int, bool, error) {

	iterator, err := keys.Iter()

	if err != nil {

		return 0, false, err
	}
	defer iterator.Close(&err)

	currentKey, err := iterator.Next()

	if err != nil {

		return 0, false, err
	}
	for currentKey != nil {

		b, err := btree.getCompareValue((*currentKey).Id())

		compareValue := sortValue.Compare(*b)

		if compareValue == -1 { // new < existing

			// new key points to the current compared key as the next key

			return iterator.Current(), false, err

		} else if compareValue == 1 { // new > existing

			// go to next key, unless this is last key

			currentKey, err = iterator.Next()

			if err != nil {

				return 0, false, err
			}

		} else { // new == existing

			return iterator.Current(), true, err
		}
	}
	return int(keys.Size()), false, err
}

// func (btree BTree[T, U]) IndexForKey(keys list.FixedSizeLinkedList[T], objectId uuid.UUID) (int, bool, error) {

// 	iterator, err := keys.Iter()

// 	if err != nil {

// 		return 0, false, err
// 	}
// 	defer iterator.Close(&err)

// 	key, err := iterator.Next()

// 	if err != nil {

// 		return 0, false, err
// 	}
// 	for key != nil {

// 		uuidCompare := btree.CompareUUIDs(objectId, (*key).Id())

// if uuidCompare == -1 { // new < existing

// 	// new key points to the current compared key as the next key

// 	return iterator.Current(), false, err

// } else if uuidCompare == 1 { // new > existing

// 	// go to next key, unless this is last key

// 	key, err = iterator.Next()

// 	if err != nil {

// 		return 0, false, err
// 	}

// } else { // new == existing

// 	return iterator.Current(), true, err
// }
// 	}
// 	return int(keys.Size()), false, err
// }

func (btree *BTree[T, U]) IndexAsChildNode(node BTreeNode) (int, error) {

	parent, exists := btree.parentage[node]

	if !exists {

		return 0, fmt.Errorf("node has no parent")
	}
	siblings := list.From[BTreeNode](btree.fileManager, parent.Children)

	for i := 0; i < siblings.Size(); i++ {

		sibling, err := siblings.Index(i)

		if err != nil {

			return 0, err
		}
		if sibling == node {

			return i, nil
		}
	}
	return 0, fmt.Errorf("node not found in parent's children")
}

func (btree *BTree[T, U]) RightMostKey(node BTreeNode) (T, BTreeNode, int, error) {

	var t T

	keys := list.From[T](btree.fileManager, node.Keys)

	// If leaf, return the rightmost key
	if node.Children == uuid.Nil {

		lastIndex := keys.Size() - 1

		key, err := keys.Index(lastIndex)

		if err != nil {

			return t, BTreeNode{}, 0, err
		}
		return key, node, lastIndex, nil
	}
	// Otherwise recurse into the rightmost child
	children := list.From[BTreeNode](btree.fileManager, node.Children)

	lastChild, err := children.Index(children.Size() - 1)

	if err != nil {

		return t, BTreeNode{}, 0, err
	}
	return btree.RightMostKey(lastChild)
}

func (btree *BTree[T, U]) MinKeys() int {

	return (btree.order / 2) - 1
}

func (btree *BTree[T, U]) ToString() string {

	str := ""

	btree.BuildString(&str, btree.root, 0)

	return str
}

func (btree *BTree[T, U]) BuildString(str *string, node BTreeNode, level int) error {

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

func (btree BTree[T, U]) CompareUUIDs(a, b uuid.UUID) int {

	return bytes.Compare(a[:], b[:]) // Lexicographic comparison
}
