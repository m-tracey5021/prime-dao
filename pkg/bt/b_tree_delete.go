package bt

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/list"
)

// https://www.programiz.com/dsa/deletion-from-a-b-tree

func (btree *BTree[T]) Delete(object T) error {

	node, nodeIndex, key, keyIndex, err := btree.Search(object)

	if err != nil {

		return err
	}
	keys := list.From[BTreeNodeKey](btree.fileManager, node.Keys)

	// if its a leaf node, then just delete the key, if the node has no
	// children after deletion, then need to handle underflow
	if node.Children == uuid.Nil {

		if key.Size() > 1 {

			return keys.Remove(keyIndex)

		} else {

			if err := keys.Remove(0); err != nil {

				return err
			}
			parent := btree.parentage[*node]

			parentKeys := list.From[BTreeNodeKey](btree.fileManager, parent.Keys)

			siblings := list.From[BTreeNode](btree.fileManager, parent.Children)

			// make sure to adjust parentage if necessary within Rotate
			rotatedLeft, err := btree.Rotate(nodeIndex, keys, parentKeys, siblings, true)

			if err != nil {

				return err
			}
			if !rotatedLeft {

				rotatedRight, err := btree.Rotate(nodeIndex, keys, parentKeys, siblings, false)

				if err != nil {

					return err
				}
				if !rotatedLeft && !rotatedRight {

					// merge case
				}
			}
			// if size of keys is less than minimum required (m/2 - 1),
			// then borrow key from immediate sibling from left to right

			// if the left sibling has extra keys, send the largest to the parent, and the parent down
			// to replace the deleted key in the child

			// if the right sibling has extra keys, send it to replace the deleted key

			// if neither has enough keys, merge them together through the parent

			return nil
		}

	} else {

		if keys.Size() > 1 {

			// if the interal node has more than the minimum number of keys, choose either its inorder predecessor or successor and move it up to the node
			// if the nodes which contain the predecessor and successor both have the minimum number of keys, merge
			// the children together either side of where the deletion took place

			// get the next key after this one wherever it is in the children
			children := list.From[BTreeNode](btree.fileManager, node.Children)

			// merge all children after this point and between next node which still remains in parent
			child, err := children.Index(keyIndex)

			if err != nil {

				return err
			}
			secondKeyIndex := keyIndex + 1

			secondChild, err := children.Index(secondKeyIndex)

			if err != nil {

				return err
			}
			childKeys := list.From[BTreeNodeKey](btree.fileManager, child.Keys)

			secondChildKeys := list.From[BTreeNodeKey](btree.fileManager, secondChild.Keys)

			if err := childKeys.AppendAll(secondChildKeys); err != nil {

				return err
			}
			if err := children.Remove(secondKeyIndex); err != nil {

				return err
			}
			return keys.Remove(keyIndex)

		} else {

			// if the internal node only has the minimum number of keys
			// delete then key and borrow from the child keys where possible, if
			// not possible, then merge the children, move the children into the empty parent
		}

		// handle underflow
		return nil
	}
}

func (btree *BTree[T]) Rotate(nodeIndex int, nodeKeys list.LinkedList[BTreeNodeKey], parentKeys list.LinkedList[BTreeNodeKey], siblings list.LinkedList[BTreeNode], rotateLeft bool) (bool, error) {

	canRotate := false

	siblingIndex := 0

	if rotateLeft {

		canRotate = nodeIndex != 0

		siblingIndex = nodeIndex - 1

	} else {

		canRotate = nodeIndex < siblings.Size()-1

		siblingIndex = nodeIndex + 1
	}
	if canRotate {

		relevantSibling, err := siblings.Index(siblingIndex)

		if err != nil {

			return false, err
		}
		relevantKeys := list.From[BTreeNodeKey](btree.fileManager, relevantSibling.Keys)

		if relevantKeys.Size() > 1 {

			keyIndex := 0

			if !rotateLeft {

				keyIndex = relevantKeys.Size() - 1
			}
			keyToTransfer, err := relevantKeys.Index(keyIndex)

			if err != nil {

				return false, err
			}
			parentKey, err := parentKeys.Index(nodeIndex)

			if err != nil {

				return false, err
			}
			if err := parentKeys.Replace(keyToTransfer, nodeIndex); err != nil {

				return false, err
			}
			if err := nodeKeys.Append(parentKey); err != nil {

				return false, err
			}
			return true, nil
		}
		return false, nil
	}
	return false, nil
}
