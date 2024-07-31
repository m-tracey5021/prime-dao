package dao

import (
	"errors"
	"io"
)

func (dao *FixedSizeDao[T]) Delete(id uint64) error {

	mainTable, err := dao.fileContainer.MainTable()

	defer func() {

		if innerErr := dao.fileContainer.Close(mainTable); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	if err != nil {

		return err
	}
	position := dao.hash(id)

	for {

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

				if err := dao.objectIO.Zero(mainTable); err != nil {

					return err
				}
				return nil

			} else {

				return dao.DeleteForCollision(id, bucketHeader.collisionTableId)
			}

		} else {

			return errors.New("object does not exist to delete")
		}
	}
}

func (dao *FixedSizeDao[T]) DeleteForCollision(id uint64, collisionTableId uint64) error {

	collisionTable, err := dao.fileContainer.CollisionTable(collisionTableId)

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
			if err := dao.objectIO.Zero(collisionTable); err != nil {

				return err
			}
			return nil
		}
	}
}
