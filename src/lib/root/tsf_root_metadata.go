package root

import (
	"reflect"
	"tsf-dao/src/lib/data/dao"
)

type TSFRootMetadata struct {
	daoIdStore dao.DaoIdStore

	daoTypeMap map[uint64]reflect.Type

	componentDaoMap map[string][]uint64 // Turn this into name of component to the corresponding daos it needs
}

func NewMetadata() TSFRootMetadata {

	return TSFRootMetadata{
		daoIdStore: dao.DaoIdStore{},

		daoTypeMap: make(map[uint64]reflect.Type),

		componentDaoMap: make(map[string][]uint64),
	}
}
