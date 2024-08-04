package dao

import (
	"errors"
	"io"
)

// func (dao *TSFHashTable[T]) DeleteObjectId(id uint64) error {

// 	initialSize := len(dao.identifierCache.ObjectIds)

// 	matches := func(element uint64) bool {

// 		return element == id
// 	}
// 	dao.identifierCache.ObjectIds = slices.DeleteFunc(dao.identifierCache.ObjectIds, matches)

// 	sizeAfterDelete := len(dao.identifierCache.ObjectIds)

// 	if sizeAfterDelete == initialSize {

// 		return errors.New("id does not exist to delete")
// 	}
// 	managingFile, err := dao.fileContainer.File(HashTableManagingFile, dao.id)

// 	if err != nil {

// 		return err
// 	}
// 	defer func() {

// 		if innerErr := dao.fileContainer.Close(managingFile); innerErr != nil {

// 			if err != nil {

// 				err = errors.Join(innerErr, err)

// 			} else {

// 				err = innerErr
// 			}
// 		}
// 	}()

// 	_, err = dao.cacheIO.WriteSizePrefixed(managingFile, dao.identifierCache)

// 	return err
// }

func (dao *TSFHashTable[T]) Delete(id uint64) error {

	mainTable, err := dao.fileContainer.Open(HashTableMainTable, dao.id)

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := dao.fileContainer.Close(mainTable); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	position := dao.hash(id)

	if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

		return err
	}
	bucketHeader, err := dao.bucketHeaderIO.Read(mainTable)

	if err != nil {

		if errors.Is(err, io.EOF) {

			return errors.New("object does not exist to delete")
		}
		return err
	}
	if bucketHeader.occupied {

		readObject, err := dao.objectIO.Read(mainTable)

		if err != nil {

			return err
		}
		if id == readObject.Id() {

			if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

				return err
			}
			bucketHeader.occupied = false

			if err := dao.bucketHeaderIO.Write(mainTable, bucketHeader); err != nil {

				return err
			}
			if err := dao.objectIO.Zero(mainTable); err != nil {

				return err
			}

		} else {

			if err := dao.DeleteForCollision(id, bucketHeader.collisionTableId); err != nil {

				return err
			}
		}
		// if err := dao.DeleteObjectId(id); err != nil {

		// 	return err
		// }
		return err

	} else {

		return errors.New("object does not exist to delete")
	}
}

func (dao *TSFHashTable[T]) DeleteForCollision(id uint64, collisionTableId uint64) error {

	collisionTable, err := dao.fileContainer.Open(HashTableCollisionTable, collisionTableId)

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := dao.fileContainer.Close(collisionTable); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	for {

		position, err := dao.fileContainer.CurrentPosition(collisionTable)

		if err != nil {

			return err
		}
		readObject, err := dao.objectIO.Read(collisionTable)

		if err != nil {

			if errors.Is(err, io.EOF) {

				return errors.New("object does not exist to delete")
			}
			return err
		}
		if id == readObject.Id() {

			if err := dao.fileContainer.GoTo(position, collisionTable); err != nil {

				return err
			}
			if err := dao.objectIO.Delete(collisionTable); err != nil {

				return err
			}
			return err
		}
	}
}
