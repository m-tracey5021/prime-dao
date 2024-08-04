package dao

import (
	"errors"
	"slices"
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

	fileManager IFileManager

	managingInfo DaoManagingInfo

	infoIO daoio.IDaoIO[DaoManagingInfo]

	indexHashTable ITSFHashTable[DaoIndex]

	metadataHashTable ITSFHashTable[DaoMetadata]

	objectIO daoio.IDaoIO[DaoIdentifiable[T]]
}

func NewDao[T any](
	id uint64,

	fileManager IFileManager,

	infoIO daoio.IDaoIO[DaoManagingInfo],

	indexHashTable ITSFHashTable[DaoIndex],

	metadataHashTable ITSFHashTable[DaoMetadata],

	objectIO daoio.IDaoIO[DaoIdentifiable[T]],

) (TSFDao[T], error) {

	managingFile, err := fileManager.Open(DaoManagingFile, id)

	defer func() {

		if innerErr := fileManager.Close(managingFile); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

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

func (dao TSFDao[T]) Get(id uint64) (*DaoIdentifiable[T], error) {

	index, err := dao.indexHashTable.Get(id)

	if err != nil {

		return nil, err
	}
	table, err := dao.fileManager.Open(DaoTable, index.id)

	if err != nil {

		return nil, err
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

		return nil, err
	}
	object, err := dao.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return nil, err
	}
	return object, err
}

func (dao *TSFDao[T]) NewObjectId() uint64 {

	return smallestMissing(dao.managingInfo.ObjectIds)
}

func (dao *TSFDao[T]) NewTableId() uint64 {

	return smallestMissing(dao.managingInfo.TableIds)
}

// TODO check pointers are correct
func (dao TSFDao[T]) AddToCache(id uint64, cache *[]uint64) {

	*cache = append(*cache, id)
}

func (dao TSFDao[T]) RemoveFromCache(id uint64, cache *[]uint64) {

	matches := func(element uint64) bool {

		return element == id
	}
	*cache = slices.DeleteFunc(*cache, matches)

}

func (dao *TSFDao[T]) AlterCache(id uint64, alteration CacheAlteration) {

	if alteration == AddTable {

		dao.AddToCache(id, &dao.managingInfo.TableIds)

	} else if alteration == AddObject {

		dao.AddToCache(id, &dao.managingInfo.ObjectIds)

	} else if alteration == RemoveTable {

		dao.RemoveFromCache(id, &dao.managingInfo.TableIds)

	} else if alteration == RemoveObject {

		dao.RemoveFromCache(id, &dao.managingInfo.ObjectIds)

	} else {

		panic("cache type not recognised")
	}
}

func (dao *TSFDao[T]) SaveCacheAlteration(id uint64, alteration CacheAlteration) error {

	managingFile, err := dao.fileManager.Open(DaoManagingFile, dao.id)

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := dao.fileManager.Close(managingFile); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	dao.AlterCache(id, alteration)

	_, err = dao.infoIO.WriteSizePrefixed(managingFile, dao.managingInfo)

	if err != nil {

		return err
	}
	return err
}

func (dao TSFDao[T]) Save(object T) (int, error) {

	metadata, err := dao.metadataHashTable.Get(dao.managingInfo.AvailableTable)

	if err != nil {

		return 0, err
	}
	if metadata.objectsWritten == dao.managingInfo.MaxObjects {

		tableId := dao.NewTableId()

		if err := dao.SaveCacheAlteration(tableId, AddTable); err != nil {

			return 0, err
		}
		dao.managingInfo.AvailableTable = tableId
	}
	table, err := dao.fileManager.Open(DaoTable, dao.managingInfo.AvailableTable)

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

	objectId := dao.NewObjectId()

	identifiable := DaoIdentifiable[T]{objectId, object}

	position, err := dao.fileManager.Size(table)

	if err := dao.fileManager.GoTo(position, table); err != nil {

		return 0, err
	}
	size, err := dao.objectIO.WriteSizePrefixed(table, identifiable)

	if err != nil {

		return 0, err
	}
	if err := dao.SaveCacheAlteration(objectId, AddObject); err != nil {

		return 0, err
	}
	index := DaoIndex{objectId, dao.managingInfo.AvailableTable, uint64(position)}

	if err := dao.indexHashTable.Save(index); err != nil {

		return 0, err
	}
	metadata.objectsWritten += 1

	if err := dao.metadataHashTable.Update(*metadata); err != nil {

		return 0, err
	}
	return size, err
}

func (dao TSFDao[T]) Update(object DaoIdentifiable[T]) (int, error) {

	index, err := dao.indexHashTable.Get(object.Id)

	if err != nil {

		return 0, err
	}
	table, err := dao.fileManager.Open(DaoTable, index.id)

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
	return dao.objectIO.Update(table, object)

	// TODO update indexes!!!
}
