package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

// func (dao TSFDao[T]) Save(object T) (*T, int, error) {

// 	table, err := dao.fileManager.OpenAndLock(fm.DaoTable, dao.tableManager.AvailableTable())

// 	defer dao.fileManager.CloseAndUnlock(table, &err)

// 	if err != nil {

// 		return nil, 0, err
// 	}
// 	object = dao.tableManager.AssignId(object)

// 	position, err := dao.fileManager.Size(table)

// 	if err := dao.fileManager.GoTo(position, table); err != nil {

// 		return nil, 0, err
// 	}
// 	size, err := dao.objectIO.WriteSizePrefixed(table, object)

// 	if err != nil {

// 		return nil, 0, err
// 	}
// 	dao.tableManager.AlterCache(object.Id(), AddObject)

// 	index := DaoIndex{object.Id(), dao.tableManager.AvailableTable(), uint64(position)}

// 	if err := dao.indexHashTable.Save(index); err != nil {

// 		return nil, 0, err
// 	}
// 	if err := dao.UpdateMetadataForSave(dao.tableManager.AvailableTable()); err != nil {

// 		return nil, 0, err
// 	}
// 	if err := dao.tableManager.SaveManagingInfo(); err != nil {

// 		return nil, 0, err
// 	}
// 	return &object, size, err
// }

type SaveProcessor[T schema.Identifiable] struct {
	object T

	fileManager fm.IFileManager

	metadataManager *DaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

type ProcessSaveResult[T schema.Identifiable] struct {
	object *T

	sizeWritten int

	err error
}

func (result *ProcessSaveResult[T]) Object() *T {

	return result.object
}

func (result *ProcessSaveResult[T]) Error() error {

	return result.err
}

func (processor SaveProcessor[T]) Process() queue.IResult[T] {

	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, processor.metadataManager.AvailableTable())

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &ProcessSaveResult[T]{nil, 0, err}
	}
	object := processor.metadataManager.AssignId(processor.object)

	position, err := processor.fileManager.Size(table)

	if err := processor.fileManager.GoTo(position, table); err != nil {

		return &ProcessSaveResult[T]{nil, 0, err}
	}
	size, err := processor.objectIO.WriteSizePrefixed(table, object)

	if err != nil {

		return &ProcessSaveResult[T]{nil, 0, err}
	}
	processor.metadataManager.UpdateForSave(table, object, uint64(position))

	return &ProcessSaveResult[T]{&object, size, err}
}
