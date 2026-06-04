package bt

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/list"
)

func (btree BTree[T]) Search(object T) (*BTreeNode, int, *BTreeNodeKey, int, error) {

	return btree.SearchRecurse(object, btree.root, 0)
}

func (btree BTree[T]) SearchRecurse(object T, node BTreeNode, nodeIndex int) (*BTreeNode, int, *BTreeNodeKey, int, error) {

	keys := list.From[BTreeNodeKey](btree.fileManager, node.Keys)

	indexForKey, found, err := btree.IndexForKey(keys, object.Id())

	if err != nil {

		return nil, 0, nil, 0, err
	}
	if found {

		key, err := keys.Index(indexForKey)

		if err != nil {

			return nil, 0, nil, 0, err
		}
		return &node, nodeIndex, &key, indexForKey, nil

	} else {

		if node.Children != uuid.Nil {

			children := list.From[BTreeNode](btree.fileManager, node.Children)

			child, err := children.Index(indexForKey)

			if err != nil {

				return nil, 0, nil, 0, err
			}
			btree.parentage[child] = node

			return btree.SearchRecurse(object, child, indexForKey)
		}
		return nil, 0, nil, 0, nil
	}
}
