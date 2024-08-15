package dao

import (
	"errors"
	"io"
	"os"
	"transformer/src/lib"
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/ht"
	"transformer/src/lib/daoio"
)

type CacheAlteration int

const (
	AddTable = 0

	AddObject

	RemoveTable

	RemoveObject
)

type TSFDao[T any] struct {
	id uint64

	fileManager fm.IFileManager

	managingInfo DaoManagingInfo

	infoIO daoio.IDaoIO[DaoManagingInfo]

	indexHashTable ht.ITSFHashTable[DaoIndex]

	metadataHashTable ht.ITSFHashTable[DaoMetadata]

	objectIO daoio.IDaoIO[DaoIdentifiable[T]]
}

func Initialise[T any](
	id uint64,

	fileManager fm.IFileManager,

	infoIO daoio.IDaoIO[DaoManagingInfo],

	indexHashTable ht.ITSFHashTable[DaoIndex],

	metadataHashTable ht.ITSFHashTable[DaoMetadata],

	objectIO daoio.IDaoIO[DaoIdentifiable[T]],

) (TSFDao[T], error) {

	managingFile, err := fileManager.OpenAndLock(fm.DaoManagingFile, id)

	defer fileManager.CloseAndUnlock(managingFile, &err)

	if err != nil {

		return TSFDao[T]{}, err
	}
	size, err := fileManager.Size(managingFile)

	if err != nil {

		return TSFDao[T]{}, err
	}
	var identifierCache *DaoManagingInfo

	if size > 0 {

		if identifierCache, err = infoIO.ReadSizePrefixed(managingFile); err != nil {

			return TSFDao[T]{}, err
		}

	} else {

		objectIds := make([]uint64, 0)

		fileIds := make([]uint64, 0)

		maxObjects := uint64(10) // Get from config

		available := uint64(0)

		identifierCache = &DaoManagingInfo{objectIds, fileIds, maxObjects, available}

		if _, err := infoIO.WriteSizePrefixed(managingFile, *identifierCache); err != nil {

			return TSFDao[T]{}, err
		}
	}
	return TSFDao[T]{id, fileManager, *identifierCache, infoIO, indexHashTable, metadataHashTable, objectIO}, err
}

func (dao *TSFDao[T]) NewObjectId() uint64 {

	return lib.NewId(dao.managingInfo.ObjectIds)
}

func (dao *TSFDao[T]) NewTableId() uint64 {

	return lib.NewId(dao.managingInfo.TableIds)
}

func (dao *TSFDao[T]) AddObjectIdToCache(id uint64) {

	dao.managingInfo.ObjectIds = append(dao.managingInfo.ObjectIds, id)
}

func (dao *TSFDao[T]) AddTableIdToCache(id uint64) {

	dao.managingInfo.TableIds = append(dao.managingInfo.TableIds, id)
}

func (dao *TSFDao[T]) RemoveObjectIdFromCache(id uint64) {

	dao.managingInfo.ObjectIds = lib.RemoveId(id, dao.managingInfo.ObjectIds)
}

func (dao *TSFDao[T]) RemoveTableIdFromCache(id uint64) {

	dao.managingInfo.TableIds = lib.RemoveId(id, dao.managingInfo.TableIds)
}

func (dao *TSFDao[T]) SaveCacheAlteration(id uint64, alteration func(id uint64)) error {

	managingFile, err := dao.fileManager.OpenAndLock(fm.DaoManagingFile, dao.id)

	defer dao.fileManager.CloseAndUnlock(managingFile, &err)

	if err != nil {

		return err
	}
	alteration(id)

	_, err = dao.infoIO.WriteSizePrefixed(managingFile, dao.managingInfo)

	if err != nil {

		return err
	}
	return err
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

func (dao *TSFDao[T]) UpdateMetadataForSave(fileIdSavedTo uint64) error {

	metadata, err := dao.metadataHashTable.Get(fileIdSavedTo)

	if err != nil {

		return err
	}
	objectsWrittenAfterSave := metadata.objectsWritten + 1

	if objectsWrittenAfterSave == dao.managingInfo.MaxObjects {

		tableId := dao.NewTableId()

		if err := dao.SaveCacheAlteration(tableId, dao.AddTableIdToCache); err != nil {

			return err
		}
		dao.managingInfo.AvailableTable = tableId
	}
	metadata.objectsWritten += 1

	if err := dao.metadataHashTable.Update(*metadata); err != nil {

		return err
	}
	return err
}

func (dao TSFDao[T]) UpdateMetadataForDeletion(table *os.File, fileIdDeletedFrom uint64) error {

	metadata, err := dao.metadataHashTable.Get(fileIdDeletedFrom)

	if err != nil {

		return err
	}
	if metadata.objectsWritten == 1 {

		if err := dao.fileManager.Remove(table); err != nil {

			return err
		}
		dao.metadataHashTable.Delete(metadata.id)

		if err := dao.SaveCacheAlteration(fileIdDeletedFrom, dao.RemoveTableIdFromCache); err != nil {

			return err
		}

	} else {

		metadata.objectsWritten -= 1

		dao.metadataHashTable.Update(*metadata)

		dao.managingInfo.AvailableTable = fileIdDeletedFrom
	}
	return err
}
