package root

import (
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/data/dao"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type DaoRoot struct {
	fileManager fm.IFileManager
}

func Connect(path string) (DaoRoot, error) {

	var exists bool

	info, err := os.Stat(path)

	if err != nil {

		exists = !os.IsNotExist(err)

	} else {

		exists = info.IsDir()
	}
	if !exists {

		return DaoRoot{}, CorruptedInstallation
	}
	fileManager := fm.NewFileManager(path, "prm")

	return DaoRoot{fileManager}, nil
}

func Dao[T schema.DescribedIdentifiable](root *DaoRoot) (*dao.Dao[T], error) {

	dao, err := dao.From[T](root.fileManager)

	if err != nil {

		return nil, err
	}
	return dao, nil
}
