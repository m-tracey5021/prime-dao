package root

import (
	"prime-dao/pkg/data/dao"
)

type TSFRootMetadata struct {
	DaoIdStore dao.DaoIdStore

	ComponentDaoMap map[string][]uint64 // Turn this into name of component to the corresponding daos it needs
}

func NewMetadata() TSFRootMetadata {

	return TSFRootMetadata{
		DaoIdStore: dao.DaoIdStore{},

		ComponentDaoMap: make(map[string][]uint64),
	}
}
