package root

import (
	"os"
	"sync"

	"github.com/m-tracey5021/prime-dao/pkg/data/dao"
	"github.com/m-tracey5021/prime-dao/pkg/data/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type DaoInit[T schema.Identifiable] interface {
	Init() dao.TSFDao[T]
}

type TSFRoot struct {
	fileManager fm.IFileManager

	metadata TSFRootMetadata

	mu sync.Mutex
}

func Connect(path string) (*TSFRoot, error) {

	var exists bool

	info, err := os.Stat(path)

	if err != nil {

		exists = !os.IsNotExist(err)

	} else {

		exists = info.IsDir()
	}
	if !exists {

		return nil, CorruptedInstallation
	}
	fileManager := fm.NewFileManager(path, "tsf")

	rootTable, err := fileManager.OpenAndLock(fm.Root)

	defer fileManager.CloseAndUnlock(rootTable, &err)

	if err != nil {

		return nil, err
	}
	size, err := fileManager.Size(rootTable)

	if size > 0 {

		metadataIO := dataio.DataIO[TSFRootMetadata]{}

		metadata, err := metadataIO.ReadSizePrefixed(rootTable)

		if err != nil {

			return nil, err
		}
		return &TSFRoot{fileManager: fileManager, metadata: metadata}, err

	} else {

		return &TSFRoot{fileManager: fileManager, metadata: NewMetadata()}, err
	}
}

func NewDaoOfType[T schema.DescribedIdentifiable](root *TSFRoot) (uint64, *dao.TSFDao[T], error) {

	id := root.metadata.DaoIdStore.NewId(&root.mu)

	dao, err := dao.From[T](root.fileManager, id)

	if err != nil {

		return 0, nil, err
	}
	return id, dao, nil
}

func GetDaoOfType[T schema.DescribedIdentifiable](root *TSFRoot, id uint64) (*dao.TSFDao[T], error) {

	return dao.From[T](root.fileManager, id)
}

func (root *TSFRoot) AssociateDaoWithComponent(componentName string, daoId ...uint64) {

	root.metadata.ComponentDaoMap[componentName] = append(root.metadata.ComponentDaoMap[componentName], daoId...)
}

func (root *TSFRoot) GetDaoForComponentName(componentName string) ([]uint64, error) {

	ids, ok := root.metadata.ComponentDaoMap[componentName]

	if ok {

		return ids, nil
	}
	return []uint64{}, ComponentDoesNotExist
}

func (root *TSFRoot) Quit() error {

	metadataIO := dataio.DataIO[TSFRootMetadata]{}

	rootTable, err := root.fileManager.OpenAndLock(fm.Root)

	defer root.fileManager.CloseAndUnlock(rootTable, &err)

	if err != nil {

		return err
	}
	if _, err := metadataIO.WriteSizePrefixed(rootTable, root.metadata); err != nil {

		return err
	}
	return nil
}
