package dao

import (
	"fmt"
	. "transformer/src/lib/dao/schema"
	daoIO "transformer/src/lib/dao_io"
)

type DaoManager[T Identifiable] struct {
	path string

	managingFile string

	managingInfo DaoManagerInfo
}

func (manager *DaoManager[T]) NewDaoManager(path, descriptor string) (*DaoManager[T], error) {

	managingFile := fmt.Sprintf("%v/%v_dao", path, descriptor)

	// file, err := os.OpenFile(managingFile, os.O_RDWR|os.O_CREATE, 0666)

	daoIO, err := daoIO.NewDaoIO[DaoManagerInfo](managingFile)

	if err != nil {

		return nil, err
	}
	size, err := daoIO.Size()

	if err != nil {

		return nil, err
	}
	var managingInfo *DaoManagerInfo

	if size > 0 {

		if managingInfo, err = daoIO.Read(); err != nil {

			return nil, err
		}

	} else {

		fileNames := make(map[uint64]string)

		available := uint64(0)

		managingInfo = &DaoManagerInfo{fileNames, available}
	}
	return &DaoManager[T]{path, managingFile, *managingInfo}, nil
}
