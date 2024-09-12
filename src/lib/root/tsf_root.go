package root

import (
	"os"
	"sync"
	"tsf-dao/src/lib/data/dao"
	"tsf-dao/src/lib/data/dataio"
	"tsf-dao/src/lib/data/fm"
	"tsf-dao/src/lib/data/schema"
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
	return &TSFRoot{
			fileManager: fm.NewFileManager(path, ""), // TODO ensure this param is useful or create new method

			metadata: NewMetadata(),
		},

		nil
}

func NewDaoOfType[T schema.DescribedIdentifiable](root *TSFRoot) (uint64, *dao.TSFDao[T], error) {

	id := root.metadata.daoIdStore.NewId(&root.mu)

	dao, err := dao.From[T](root.fileManager, id)

	if err != nil {

		return 0, nil, err
	}
	return id, dao, nil
}

func GetDaoOfType[T schema.DescribedIdentifiable](root *TSFRoot, id uint64) (*dao.TSFDao[T], error) {

	return dao.From[T](root.fileManager, id)
}

func (root *TSFRoot) NewComponentForName(componentName string) {

	root.metadata.componentDaoMap[componentName] = make([]uint64, 0)
}

func (root *TSFRoot) AssociateDaoWithComponent(daoId uint64, componentName string) {

	root.metadata.componentDaoMap[componentName] = append(root.metadata.componentDaoMap[componentName], daoId)
}

func (root *TSFRoot) GetDaoForComponentName(componentName string) ([]uint64, error) {

	ids, ok := root.metadata.componentDaoMap[componentName]

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
