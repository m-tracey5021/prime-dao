package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

// func (dao TSFDao[T]) Get(id uint64) (*T, error) {

// 	index, err := dao.indexHashTable.Get(id)

// 	if err != nil {

// 		return nil, err
// 	}
// 	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, index.id)

// 	defer dao.fileManager.CloseAndUnlock(table, &err)

// 	if err != nil {

// 		return nil, err
// 	}
// 	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

// 		return nil, err
// 	}
// 	object, err := dao.objectIO.ReadSizePrefixed(table)

// 	if err != nil {

// 		return object, err
// 	}
// 	return nil, err
// }

func (dao TSFDao[T]) ProcessGet(id uint64) (*T, error) {

	getRequest := dao.requstFactory.CreateGetRequest(dao.id, dao.fileManager, dao.metadataManager, id)

	result := dao.queue.ProcessSync(getRequest)

	return result.Object(), result.Error()
}

func (dao TSFDao[T]) ProcessGetAsync(ids ...uint64) {

	requests := make([]queue.IProcessor[T], 0)

	for _, id := range ids {

		getRequest := dao.requstFactory.CreateGetRequest(dao.id, dao.fileManager, dao.metadataManager, id)

		requests = append(requests, getRequest)
	}
	dao.queue.ProcessAsync(requests...)
}

type GetProcessor[T schema.Identifiable] struct {
	id uint64

	fileManager fm.IFileManager

	metadataManager *DaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

type ProcessGetResult[T schema.Identifiable] struct {
	object *T

	err error
}

func (result *ProcessGetResult[T]) Object() *T {

	return result.object
}

func (result *ProcessGetResult[T]) Error() error {

	return result.err
}

func (processor GetProcessor[T]) Process() queue.IResult[T] {

	index, err := processor.metadataManager.GetIndex(processor.id)

	if err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.id)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	object, err := processor.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return &ProcessGetResult[T]{object, err}
	}
	return &ProcessGetResult[T]{nil, err}
}
