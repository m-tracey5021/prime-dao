package dao

import (
	"errors"
	"io"
)

func (dao *FixedSizeDao[T]) Update(object T) error {

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

			return errors.New("object does not exist to update")
		}
		return err
	}
	if bucketHeader.occupied {

		readObject, err := dao.objectIO.Read(mainTable)

		if err != nil {

			return err
		}
		if object.Id() == readObject.Id() {

			if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

				return err
			}
			if err := dao.objectIO.Write(mainTable, object); err != nil {

				return err
			}
			return err

		} else {

			return dao.UpdateForCollision(object, bucketHeader.collisionTableId)
		}

	} else {

		return errors.New("object does not exist to update")
	}
}

func (dao *FixedSizeDao[T]) UpdateForCollision(object T, collisionTableId uint64) error {

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

		position, err := dao.fileContainer.CurrentPosition(collisionTable)

		if err != nil {

			return err
		}
		readObject, err := dao.objectIO.Read(collisionTable)

		if err != nil {

			if errors.Is(err, io.EOF) {

				return errors.New("object does not exist to update")
			}
			return err
		}
		if object.Id() == readObject.Id() {

			if err := dao.fileContainer.GoTo(position, collisionTable); err != nil {

				return err
			}
			if err := dao.objectIO.Write(collisionTable, object); err != nil {

				return err
			}
			return err
		}
	}
}
