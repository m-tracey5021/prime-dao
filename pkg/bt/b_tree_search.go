package bt

import (
	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/list"
)

/*
Figure out what i need to pass in here to
compare on whatever value is selected to be the relevant
comparison, maybe need to introduce another type parameter
*/
func (btree BTree[T]) Search(object T) (*BTreeNode, int, *T, int, error) {

	return btree.SearchRecurse(object, btree.root, 0)
}

func (btree BTree[T]) SearchRecurse(object T, node BTreeNode, nodeIndex int) (*BTreeNode, int, *T, int, error) {

	keys := list.From[T](btree.fileManager, node.Keys)

	indexForKey, found, err := btree.IndexForKey(keys, object)

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
