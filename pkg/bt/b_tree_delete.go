package bt

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/list"
)

// https://www.programiz.com/dsa/deletion-from-a-b-tree

func (btree *BTree[T, U]) Delete(object T) error {

	node, nodeIndex, _, keyIndex, err := btree.Search(object.Id())

	if err != nil {

		return err
	}
	keys := list.From[T](btree.fileManager, node.Keys)

	// if its a leaf node, then just delete the key, if the node has no
	// children after deletion, then need to handle underflow
	if node.Children == uuid.Nil {

		if err := keys.Remove(keyIndex); err != nil {

			return err
		}
		// No underflow if node still has enough keys
		if keys.Size() >= btree.MinKeys() {

			return nil
		}
		return btree.Underflow(*node, nodeIndex)

	} else {

		// Internal node, so replace with in-order predecessor
		children := list.From[BTreeNode](btree.fileManager, node.Children)

		leftChild, err := children.Index(keyIndex)

		if err != nil {

			return err
		}
		predecessor, predecessorNode, predecessorKeyIndex, err := btree.RightMostKey(leftChild)

		if err != nil {

			return err
		}
		// Replace the deleted key with the predecessor
		if err := keys.Replace(predecessor, keyIndex); err != nil {

			return err
		}
		// Delete the predecessor from its leaf node
		predecessorKeys := list.From[BTreeNodeKey](btree.fileManager, predecessorNode.Keys)

		if err := predecessorKeys.Remove(predecessorKeyIndex); err != nil {

			return err
		}
		if predecessorKeys.Size() < btree.MinKeys() {

			predecessorNodeIndex, err := btree.IndexAsChildNode(predecessorNode)

			if err != nil {

				return err
			}
			return btree.Underflow(predecessorNode, predecessorNodeIndex)
		}
		return nil
	}
}

func (btree *BTree[T, U]) Underflow(node BTreeNode, nodeIndex int) error {

	parent, exists := btree.parentage[node]

	if !exists {
		// Node is root, underflow at root is fine (tree shrinks)
		return nil
	}
	parentKeys := list.From[BTreeNodeKey](btree.fileManager, parent.Keys)

	siblings := list.From[BTreeNode](btree.fileManager, parent.Children)

	nodeKeys := list.From[BTreeNodeKey](btree.fileManager, node.Keys)

	// Try rotating from left sibling first
	rotatedLeft, err := btree.Rotate(nodeIndex, nodeKeys, parentKeys, siblings, true)

	if err != nil {

		return err
	}
	if rotatedLeft {

		return nil
	}
	// Try rotating from right sibling
	rotatedRight, err := btree.Rotate(nodeIndex, nodeKeys, parentKeys, siblings, false)

	if err != nil {
		return err

	}
	if rotatedRight {
		return nil

	}
	// Neither rotation worked, so merge instead
	return btree.Merge(node, nodeIndex, nodeKeys, parentKeys, siblings)
}

func (btree *BTree[T, U]) Rotate(nodeIndex int, nodeKeys list.FixedSizeLinkedList[BTreeNodeKey], parentKeys list.FixedSizeLinkedList[BTreeNodeKey], siblings list.FixedSizeLinkedList[BTreeNode], rotateLeft bool) (bool, error) {

	canRotate := false

	siblingIndex := 0

	if rotateLeft {

		canRotate = nodeIndex != 0

		siblingIndex = nodeIndex - 1

	} else {

		canRotate = nodeIndex < siblings.Size()-1

		siblingIndex = nodeIndex + 1
	}
	if !canRotate {

		return false, nil
	}
	relevantSibling, err := siblings.Index(siblingIndex)

	if err != nil {

		return false, err
	}
	relevantKeys := list.From[BTreeNodeKey](btree.fileManager, relevantSibling.Keys)

	if relevantKeys.Size() <= btree.MinKeys() {

		return false, nil
	}
	// Left rotation: borrow largest from left sibling
	// Right rotation: borrow smallest from right sibling
	keyIndex := relevantKeys.Size() - 1

	if !rotateLeft {

		keyIndex = 0
	}
	keyToTransfer, err := relevantKeys.Index(keyIndex)
	if err != nil {

		return false, err
	}
	// Separator key in parent sits between the two nodes
	parentKeyIndex := nodeIndex - 1

	if !rotateLeft {

		parentKeyIndex = nodeIndex
	}
	parentKey, err := parentKeys.Index(parentKeyIndex)

	if err != nil {

		return false, err
	}
	// Remove transferred key from sibling
	if err := relevantKeys.Remove(keyIndex); err != nil {

		return false, err
	}
	// Sibling's key replaces the separator in parent
	if err := parentKeys.Replace(keyToTransfer, parentKeyIndex); err != nil {

		return false, err
	}
	// Parent's old separator key goes into the underflowing node
	if err := nodeKeys.Append(parentKey); err != nil {

		return false, err
	}
	return true, nil
}

// merge combines a node with a sibling, pulling the separator key down from the parent
func (btree *BTree[T, U]) Merge(node BTreeNode, nodeIndex int, nodeKeys list.FixedSizeLinkedList[BTreeNodeKey], parentKeys list.FixedSizeLinkedList[BTreeNodeKey], siblings list.FixedSizeLinkedList[BTreeNode]) error {

	// Prefer merging with left sibling, otherwise right
	var leftIndex, rightIndex int

	if nodeIndex > 0 {

		leftIndex = nodeIndex - 1

		rightIndex = nodeIndex

	} else {

		leftIndex = nodeIndex

		rightIndex = nodeIndex + 1
	}
	separatorIndex := leftIndex

	leftNode, err := siblings.Index(leftIndex)

	if err != nil {

		return err
	}
	rightNode, err := siblings.Index(rightIndex)

	if err != nil {

		return err
	}
	leftKeys := list.From[BTreeNodeKey](btree.fileManager, leftNode.Keys)

	rightKeys := list.From[BTreeNodeKey](btree.fileManager, rightNode.Keys)

	// Pull separator key down from parent into left node
	separatorKey, err := parentKeys.Index(separatorIndex)

	if err != nil {

		return err
	}
	if err := leftKeys.Append(separatorKey); err != nil {

		return err
	}
	// Append all right node keys into left node
	if err := leftKeys.AppendAll(rightKeys); err != nil {

		return err
	}
	// If internal nodes, move right node's children into left node
	if rightNode.Children != uuid.Nil {

		leftChildren := list.From[BTreeNode](btree.fileManager, leftNode.Children)

		rightChildren := list.From[BTreeNode](btree.fileManager, rightNode.Children)

		if err := leftChildren.AppendAll(rightChildren); err != nil {

			return err
		}
		// Update parentage for moved children
		for i := 0; i < rightChildren.Size(); i++ {

			child, err := rightChildren.Index(i)

			if err != nil {

				return err
			}
			btree.parentage[child] = leftNode
		}
	}
	// Remove the separator key and right node from parent
	if err := parentKeys.Remove(separatorIndex); err != nil {

		return err
	}
	if err := siblings.Remove(rightIndex); err != nil {

		return err
	}
	// Parent may now be underflowing, recurse upward
	parent := btree.parentage[node]

	if parentKeys.Size() < btree.MinKeys() {

		parentIndex, err := btree.IndexAsChildNode(parent)

		if err != nil {

			return err
		}
		return btree.Underflow(parent, parentIndex)
	}
	return nil
}
