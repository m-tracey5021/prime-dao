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

	indexForKey, found, err := btree.IndexForKey(keys, object)

	if err != nil {

		return err
	}
	if found {

		return errors.New("object already exists")
	}
	if node.Children == uuid.Nil { // we have found the right leaf

		if keys.Size() == btree.order-1 { // leaf is full

			return btree.SplitNode(object, parent, node, indexForKey)

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

func (btree *BTree[T, U]) SplitNode(object T, parent *BTreeNode, child BTreeNode, index int) error {

	originalKeys := list.From[T](btree.fileManager, child.Keys)

	originalKeys.Insert(object, index)

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
	newRight := BTreeNode{rightKeys.Id(), uuid.Nil}

	if parent != nil {

		parentKeys := list.From[T](btree.fileManager, parent.Keys)

		if parentKeys.Size() == btree.order-1 { // if parent is full

			// figure out how to send back up the split
			return nil

		} else {

			if middleKeyIndex == parentKeys.Size() {

				if err := parentKeys.Append(middleKey); err != nil {

					return err
				}

			} else {

				if err := parentKeys.Insert(middleKey, middleKeyIndex); err != nil {

					return err
				}
			}
			btree.parentage[newRight] = *parent

			parentChildNodes := list.From[BTreeNode](btree.fileManager, parent.Children)

			return parentChildNodes.Insert(newRight, middleKeyIndex+1)
		}

	} else {

		rootKeys := list.NewList[T](btree.fileManager)

		rootChildren := list.NewList[BTreeNode](btree.fileManager)

		newRoot := BTreeNode{rootKeys.Id(), rootChildren.Id()}

		rootKeys.Append(middleKey)

		rootChildren.Append(child)

		rootChildren.Append(newRight)

		return btree.SaveRoot(newRoot)
	}
}
