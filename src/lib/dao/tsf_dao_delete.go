package dao

import (
	"errors"
	"io"
	"os"
)

func (dao TSFDao[T]) Delete(object DaoIdentifiable[T]) (int, error) {

	index, err := dao.indexHashTable.Get(object.Id)

	if err != nil {

		return 0, err
	}
	table, err := dao.fileManager.Open(DaoTable, index.fileId)

	if err != nil {

		return 0, err
	}
	defer func() {

		if innerErr := dao.fileManager.Close(table); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	if err := dao.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return 0, err
	}
	deletedSize, err := dao.objectIO.Delete(table)

	if err != nil {

		return 0, err
	}
	if err := dao.SaveCacheAlteration(object.Id, RemoveObject); err != nil {

		return 0, err
	}
	dao.UpdateMetadataForDeletion(table, index.fileId)

	dao.UpdateIndexes(table, index.fileId, index.filePosition)

	if err := dao.indexHashTable.Delete(index.id); err != nil {

		return 0, err
	}
	return deletedSize, err
}

func (dao TSFDao[T]) UpdateIndexes(table *os.File, fileId, filePosition uint64) error {

	if err := dao.fileManager.GoTo(int(filePosition), table); err != nil {

		return err
	}
	for {

		currentPosition, err := dao.fileManager.CurrentPosition(table)

		if err != nil {

			return err
		}
		readObject, err := dao.objectIO.ReadSizePrefixed(table)

		if err != nil {

			if errors.Is(err, io.EOF) {

				break
			}
			return err
		}
		updatedIndex := DaoIndex{readObject.Id, fileId, uint64(currentPosition)}

		if err := dao.indexHashTable.Update(updatedIndex); err != nil {

			return err
		}
	}
	return nil
}

func (dao TSFDao[T]) UpdateMetadataForDeletion(table *os.File, fileId uint64) error {

	metadata, err := dao.metadataHashTable.Get(fileId)

	if err != nil {

		return err
	}
	if metadata.objectsWritten == 1 {

		if err := dao.fileManager.Remove(table); err != nil {

			return err
		}
		dao.metadataHashTable.Delete(metadata.id)

		if err := dao.SaveCacheAlteration(fileId, RemoveTable); err != nil {

			return err
		}

	} else {

		metadata.objectsWritten -= 1

		dao.metadataHashTable.Update(*metadata)

		dao.managingInfo.AvailableTable = fileId
	}
	return err
}
