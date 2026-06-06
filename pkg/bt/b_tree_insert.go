package bt

import (
	"errors"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/list"
)

func (btree *BTree[T, U]) Insert(object T) error {

	if btree.root.Keys == uuid.Nil {

		rootKeys := list.NewList[T](btree.fileManager)

		root := BTreeNode{rootKeys.Id(), uuid.Nil}

		if err := rootKeys.Append(object); err != nil {

			return err
		}
		return btree.SaveRoot(root)

	} else {

		return btree.InsertSearchRecurse(object, nil, nil, btree.root)
	}
}

func (btree *BTree[T, U]) InsertSearchRecurse(object T, parent *BTreeNode, parentIndex *int, node BTreeNode) error {

	keys := list.From[T](btree.fileManager, node.Keys)

	indexForKey, found, err := btree.IndexForNonExistingKeyOnSave(keys, object)

	if err != nil {

		return err
	}
	if found {

		return errors.New("object already exists")
	}
	if node.Children == uuid.Nil { // we have found the right leaf

		if btree.LeafIsFull(keys) { // leaf is full

			// return btree.SplitNode(object, parent, node, indexForKey)

			if err := keys.Insert(object, indexForKey); err != nil {

				return err
			}
			return btree.FixOverflow(node, parent)

		} else { // insert into child keys

			if indexForKey == keys.Size() {

				return keys.Append(object)

			} else {

				return keys.Insert(object, indexForKey)
			}
		}

	} else {

		children := list.From[BTreeNode](btree.fileManager, node.Children)

		child, err := children.Index(indexForKey)

		if err != nil {

			return err
		}
		btree.parentage[child] = node

		return btree.InsertSearchRecurse(object, &node, &indexForKey, child)
	}
}

// func (btree *BTree[T, U]) SplitNode(object T, parent *BTreeNode, indexOfParent int, child BTreeNode, indexOfObjectInChild int) error {

// 	originalKeys := list.From[T](btree.fileManager, child.Keys)

// 	originalKeys.Insert(object, indexOfObjectInChild)

// 	middleKeyIndex := originalKeys.Size() / 2

// 	middleKey, err := originalKeys.Index(int(middleKeyIndex))

// 	if err != nil {

// 		return err
// 	}
// 	rightKeys, err := originalKeys.Copy(btree.fileManager, int(middleKeyIndex)+1)

// 	if err != nil {

// 		return err
// 	}
// 	// original keys are now the left side?
// 	if err := originalKeys.Truncate(int(middleKeyIndex)); err != nil {

// 		return err
// 	}
// 	newRight := BTreeNode{rightKeys.Id(), uuid.Nil}

// 	if parent != nil {

// 		parentKeys := list.From[T](btree.fileManager, parent.Keys)

// 		if err != nil {

// 			return err
// 		}
// 		if parentKeys.Size() == btree.order-1 { // if parent is full

// 			// add first anyway, recursive call can deal with overflow

// 			// figure out how to send back up the split
// 			// get where the index of the middle key goes in the parent keys

// 			// get parent node, and grandparent node
// 			grandparent, ok := btree.parentage[*parent]

// 			if ok { // there is a grandparent

// 				return btree.SplitNode(middleKey, &grandparent, *parent, indexForMiddleKey)

// 			} else {

// 				return btree.SplitNode(middleKey, nil, *parent, indexForMiddleKey)
// 			}

// 		} else {

// 			if indexOfParent == parentKeys.Size() {

// 				if err := parentKeys.Append(middleKey); err != nil {

// 					return err
// 				}

// 			} else {

// 				if err := parentKeys.Insert(middleKey, indexOfParent); err != nil {

// 					return err
// 				}
// 			}
// 			btree.parentage[newRight] = *parent

// 			parentChildNodes := list.From[BTreeNode](btree.fileManager, parent.Children)

// 			return parentChildNodes.Insert(newRight, middleKeyIndex+1)
// 		}

// 	} else {

// 		rootKeys := list.NewList[T](btree.fileManager)

// 		rootChildren := list.NewList[BTreeNode](btree.fileManager)

// 		newRoot := BTreeNode{rootKeys.Id(), rootChildren.Id()}

// 		rootKeys.Append(middleKey)

// 		rootChildren.Append(child)

// 		rootChildren.Append(newRight)

// 		return btree.SaveRoot(newRoot)
// 	}
// }

func (btree *BTree[T, U]) FixOverflow(overFlowNode BTreeNode, parent *BTreeNode) error {

	// get middle key from OverflowNode

	// get left half of keys

	// get right half of keys

	// createa new node for each half with those keys, call this newChildNodeA and B

	// split childNodes of overflowNode by the index of the middleKey

	// attach left half of nodes to newChildNodeA and right half to B

	// attach middle key to the index it belongs to in the parent node

	originalKeys := list.From[T](btree.fileManager, overFlowNode.Keys)

	middleKeyIndex := originalKeys.Size() / 2

	middleKey, err := originalKeys.Index(int(middleKeyIndex))

	if err != nil {

		return err
	}
	rightKeys, err := originalKeys.Copy(btree.fileManager, int(middleKeyIndex)+1)

	if err != nil {

		return err
	}
	if err := originalKeys.Truncate(int(middleKeyIndex)); err != nil {

		return err
	}
	rightChildrenId := uuid.Nil

	if overFlowNode.Children != uuid.Nil {

		originalChildren := list.From[BTreeNode](btree.fileManager, overFlowNode.Children)

		rightChildren, err := originalChildren.Copy(btree.fileManager, int(middleKeyIndex)+1)

		rightChildrenId = rightChildren.Id()

		if err != nil {

			return err
		}
		if err := originalChildren.Truncate(int(middleKeyIndex)); err != nil {

			return err
		}
	}
	newRight := BTreeNode{rightKeys.Id(), rightChildrenId}

	// newRight also needs its child nodes attached if it had any

	if parent != nil {

		parentKeys := list.From[T](btree.fileManager, parent.Keys)

		indexForMiddleKeyInParent, _, err := btree.IndexForNonExistingKeyOnSave(parentKeys, middleKey)

		if err != nil {

			return err
		}
		if indexForMiddleKeyInParent == parentKeys.Size() {

			if err := parentKeys.Append(middleKey); err != nil {

				return err
			}

		} else {

			if err := parentKeys.Insert(middleKey, indexForMiddleKeyInParent); err != nil {

				return err
			}
		}
		btree.parentage[newRight] = *parent

		parentChildNodes := list.From[BTreeNode](btree.fileManager, parent.Children)

		if err := parentChildNodes.Insert(newRight, indexForMiddleKeyInParent+1); err != nil {

			return err
		}
		if btree.HasOverFlowed(parentKeys) { // if parent has already overflowed because of the above

			// insert middle key into parentKeys
			// add new right as new node on indexOfMiddleInParent +1
			grandparent, ok := btree.parentage[*parent]

			if ok { // there is a grandparent

				return btree.FixOverflow(*parent, &grandparent)

			} else {

				return btree.FixOverflow(*parent, nil)
			}
		}
		return nil

	} else {

		rootKeys := list.NewList[T](btree.fileManager)

		rootChildren := list.NewList[BTreeNode](btree.fileManager)

		newRoot := BTreeNode{rootKeys.Id(), rootChildren.Id()}

		rootKeys.Append(middleKey)

		rootChildren.Append(overFlowNode)

		rootChildren.Append(newRight)

		return btree.SaveRoot(newRoot)
	}
}

func (btree BTree[T, U]) HasOverFlowed(keys list.FixedSizeLinkedList[T]) bool {

	return keys.Size() > btree.order-1
}

func (btree BTree[T, U]) LeafIsFull(keys list.FixedSizeLinkedList[T]) bool {

	return keys.Size() == btree.order-1
}
