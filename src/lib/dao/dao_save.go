package dao

import (
	"errors"
	"io"
)

func (dao *FixedSizeDao[T]) SaveNewObjectId(id uint64) error {

	managingFile, err := dao.fileContainer.ManagingFile()

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := dao.fileContainer.Close(managingFile); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	dao.identifierCache.ObjectIds = append(dao.identifierCache.ObjectIds, id)

	_, err = dao.cacheIO.WriteSizePrefixed(managingFile, dao.identifierCache)

	if err != nil {

		return err
	}
	return err
}

func (dao *FixedSizeDao[T]) Save(object T) error {

	mainTable, err := dao.fileContainer.MainTable()

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

	position := dao.hash(object.Id())

	if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

		return err
	}
	bucketHeader, err := dao.bucketHeaderIO.Read(mainTable)

	if err != nil {

		if errors.Is(err, io.EOF) {

			err = nil

			bucketHeader = DaoBucketHeader{false, false, 0}

		} else {

			return err
		}
	}
	if bucketHeader.occupied {

		readObject, err := dao.objectIO.Read(mainTable)

		if err != nil {

			return err
		}
		if object.Id() == readObject.Id() {

			return errors.New("object already exists, cannot save new")
		}
		if !bucketHeader.previousCollision {

			collisionTableId := dao.NewCollisionTableId()

			if err := dao.SaveForCollision(object, collisionTableId); err != nil {

				return err
			}
			updatedBucketHeader := DaoBucketHeader{true, true, collisionTableId}

			if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

				return err
			}
			if err := dao.bucketHeaderIO.Write(mainTable, updatedBucketHeader); err != nil {

				return nil
			}

		} else {

			if err := dao.SaveForCollision(object, bucketHeader.collisionTableId); err != nil {

				return err
			}
		}

	} else {

		bucketHeader = DaoBucketHeader{true, false, 0}

		if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

			return err
		}
		if err := dao.bucketHeaderIO.Write(mainTable, bucketHeader); err != nil {

			return err
		}
		if err := dao.objectIO.Write(mainTable, object); err != nil {

			return err
		}
	}
	if err := dao.SaveNewObjectId(object.Id()); err != nil {

		return err
	}
	return err
}

func (dao *FixedSizeDao[T]) SaveForCollision(object T, collisionTableId uint64) error {

	collisionTable, err := dao.fileContainer.CollisionTable(collisionTableId)

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

		readObject, err := dao.objectIO.Read(collisionTable)

		if err != nil {

			if errors.Is(err, io.EOF) {

				if err := dao.objectIO.Write(collisionTable, object); err != nil {

					return err
				}
				err = nil

				return err
			}
			return err
		}
		if object.Id() == readObject.Id() {

			return errors.New("object already exists, cannot save new")
		}
	}
}
