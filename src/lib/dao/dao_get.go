package dao

import (
	"errors"
	"io"
)

func (dao *FixedSizeDao[T]) Get(id uint64) (*T, error) {

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

		return nil, err
	}
	position := dao.hash(id)

	for {

		if err := dao.fileContainer.GoTo(position, mainTable); err != nil {

			return nil, err
		}
		bucketHeader, err := dao.bucketHeaderIO.Read(mainTable)

		if err != nil {

			if errors.Is(err, io.EOF) {

				return nil, nil
			}
			return nil, err
		}
		if bucketHeader.occupied {

			object, err := dao.objectIO.Read(mainTable)

			if err != nil {

				return nil, err
			}
			if id == object.Id() {

				return &object, nil

			} else {

				return dao.GetForCollision(id, bucketHeader.collisionTableId)
			}

		} else {

			return nil, nil
		}
	}
}

func (dao *FixedSizeDao[T]) GetForCollision(id, collisionTableId uint64) (*T, error) {

	collisionTable, err := dao.fileContainer.CollisionTable(collisionTableId)

	if err != nil {

		return nil, err
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

	if err != nil {

		return nil, err
	}
	for {

		object, err := dao.objectIO.Read(collisionTable)

		if err != nil {

			if errors.Is(err, io.EOF) {

				return nil, nil
			}
			return nil, err
		}
		if id == object.Id() {

			return &object, nil
		}
	}
}
