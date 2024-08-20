package dao

import (
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type TSFDao[T schema.Identifiable] struct {
	id uint64

	fileManager fm.IFileManager

	// tableManager IDaoTableManager[T]

	// indexHashTable ht.ITSFHashTable[DaoIndex]

	// metadataHashTable ht.ITSFHashTable[DaoMetadata]

	// objectIO dataio.IDataIO[T]

	metadataManager *DaoMetadataManager[T]

	requstFactory DaoRequestFactory[T]

	queue *queue.TSFQueueB[T]
}

func New[T schema.Identifiable](path, descriptor string, id uint64) (*TSFDao[T], error) {

	fileManager := fm.NewFileManager(path, descriptor)

	// indexHashTable := ht.FromFileManager[DaoIndex](id, fileManager)

	// metadataHashTable := ht.FromFileManager[DaoMetadata](id, fileManager)

	// tableManager, err := NewDaoTableManager[T](id, fileManager)

	metadataManager, err := NewMetadataManager[T](id, fileManager)

	if err != nil {

		return nil, err
	}

	// objectIO := dataio.DataIO[T]{}

	// queue := queue.TSFQueueB[T]{}

	// requstFactory := DaoRequestFactoryB[T]{}

	return &TSFDao[T]{

			id: id,

			fileManager: fileManager,

			metadataManager: metadataManager,

			requstFactory: DaoRequestFactory[T]{},

			queue: &queue.TSFQueueB[T]{},
		},

		err
}

// func (dao TSFDao[T]) UpdateIndexes(table *os.File, fileId, filePosition uint64) error {

// 	if err := dao.fileManager.GoTo(int(filePosition), table); err != nil {

// 		return err
// 	}
// 	for {

// 		currentPosition, err := dao.fileManager.CurrentPosition(table)

// 		if err != nil {

// 			return err
// 		}
// 		readObject, err := dao.objectIO.ReadSizePrefixed(table)

// 		if err != nil {

// 			if errors.Is(err, io.EOF) {

// 				break
// 			}
// 			return err
// 		}
// 		updatedIndex := DaoIndex{(*readObject).Id(), fileId, uint64(currentPosition)}

// 		if err := dao.indexHashTable.Update(updatedIndex); err != nil {

// 			return err
// 		}
// 	}
// 	return nil
// }

// func (dao *TSFDao[T]) UpdateMetadataForSave(fileIdSavedTo uint64) error {

// 	metadata, err := dao.metadataHashTable.Get(fileIdSavedTo)

// 	if err != nil {

// 		return err
// 	}
// 	objectsWrittenAfterSave := metadata.objectsWritten + 1

// 	if objectsWrittenAfterSave == dao.tableManager.MaxObjects() {

// 		tableId := dao.tableManager.NewTableId()

// 		dao.tableManager.AlterCache(tableId, AddTable)

// 		dao.tableManager.SetAvailableTable(tableId)
// 	}
// 	metadata.objectsWritten += 1

// 	if err := dao.metadataHashTable.Update(*metadata); err != nil {

// 		return err
// 	}
// 	return err
// }

// func (dao TSFDao[T]) UpdateMetadataForDeletion(table *os.File, fileIdDeletedFrom uint64) error {

// 	metadata, err := dao.metadataHashTable.Get(fileIdDeletedFrom)

// 	if err != nil {

// 		return err
// 	}
// 	if metadata.objectsWritten == 1 {

// 		if err := dao.fileManager.Remove(table); err != nil {

// 			return err
// 		}
// 		dao.metadataHashTable.Delete(metadata.id)

// 		dao.tableManager.AlterCache(fileIdDeletedFrom, RemoveTable)

// 	} else {

// 		metadata.objectsWritten -= 1

// 		dao.metadataHashTable.Update(*metadata)

// 		dao.tableManager.SetAvailableTable(fileIdDeletedFrom)
// 	}
// 	return err
// }
