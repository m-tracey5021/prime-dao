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
func (btree BTree[T, U]) Search(id uuid.UUID) (*BTreeNode, int, *T, int, error) {

	return btree.SearchRecurse(id, btree.root, 0)
}

func (btree BTree[T, U]) SearchRecurse(id uuid.UUID, node BTreeNode, nodeIndex int) (*BTreeNode, int, *T, int, error) {

	keys := list.From[T](btree.fileManager, node.Keys)

	indexForKey, found, err := btree.IndexForExistingKey(keys, id)

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

			return btree.SearchRecurse(id, child, indexForKey)
		}
		return nil, 0, nil, 0, nil
	}
}
