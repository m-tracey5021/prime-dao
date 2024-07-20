package dao

import (
	"errors"
	"io"
	. "transformer/src/lib/dao/schema"
	. "transformer/src/lib/dao_io"
	"unsafe"
)

type HashableDao[T Hashable] struct {
	daoIOFactory IDaoIOFactory[T]

	managingFile string

	tableSize int

	bucketSize int
}

func NewHashableDao[T Hashable](daoIOFactory IDaoIOFactory[T], managingFile string, tableSize int) HashableDao[T] {

	bucketSize := int(unsafe.Sizeof(*new(T)))

	return HashableDao[T]{daoIOFactory, managingFile, tableSize, bucketSize}
}

func (dao *HashableDao[T]) io() (IDaoIO[T], error) {

	return dao.daoIOFactory.Create(dao.managingFile)
}

func (dao *HashableDao[T]) hash(id uint64) int {

	hash := (int(id)*2 + 1) % dao.tableSize

	return hash + dao.bucketSize
}

func (dao *HashableDao[T]) rehash() int {

	return dao.tableSize + dao.bucketSize
}

func (dao *HashableDao[T]) NewId() (uint64, error) {

	var id uint64 = 0

	index, err := dao.Get(id)

	if err != nil {

		return 0, err
	}
	for index != nil {

		id += 1

		index, err = dao.Get(id)

		if err != nil {

			return 0, nil
		}
	}
	return id, nil
}

func (dao *HashableDao[T]) Get(id uint64) (*T, error) {

	daoIO, err := dao.io()

	defer func() {

		if innerErr := daoIO.Close(); innerErr != nil {

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

	// Add condition to break this loop somehow
	for {

		occupied, err := daoIO.ReadBoolAt(position)

		if err != nil {

			if errors.Is(err, io.EOF) {

				return nil, nil
			}
			return nil, err
		}
		if occupied {

			read, err := daoIO.ReadAt(position)

			if err != nil {

				return nil, nil
			}
			if id == (*read).Id() {

				return read, nil

			} else {

				position += dao.rehash()
			}

		} else {

			return nil, nil
		}
	}
}

func (dao *HashableDao[T]) GetAll() ([]T, error) {

	daoIO, err := dao.io()

	defer func() {

		if innerErr := daoIO.Close(); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()
	results := make([]T, 0)

	position := 0

	for {

		occupied, err := daoIO.ReadBoolAt(position)

		if err != nil {

			return nil, err
		}
		if occupied {

			read, err := daoIO.Read()

			if err != nil {

				if errors.Is(err, io.EOF) {

					return results, nil
				}
				return nil, err
			}
			results = append(results, *read)

		} else {

			position += dao.bucketSize
		}
	}
}

func (dao *HashableDao[T]) Save(object T) error {

	daoIO, err := dao.io()

	defer func() {

		if innerErr := daoIO.Close(); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()
	position := dao.hash(object.Id())

	for {

		occupied, err := daoIO.ReadBoolAt(position)

		if err != nil {

			return err
		}
		if occupied {

			read, err := daoIO.Read()

			if err != nil {

				return err
			}
			if object.Id() == (*read).Id() {

				return errors.New("object already exists, cannot save new")

			} else {

				position += dao.rehash()
			}

		} else {

			if err := daoIO.WriteBoolAt(true, position); err != nil {

				return err
			}
			if _, err := daoIO.Write(object); err != nil {

				return err
			}
			return nil
		}
	}
}

func (dao *HashableDao[T]) Update(object T) error {

	daoIO, err := dao.io()

	defer func() {

		if innerErr := daoIO.Close(); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()
	position := dao.hash(object.Id())

	for {

		occupied, err := daoIO.ReadBoolAt(position)

		if err != nil {

			return err
		}
		if occupied {

			read, err := daoIO.Read()

			if err != nil {

				return err
			}
			if object.Id() == (*read).Id() {

				if err := daoIO.UpdateAt(object, position+1); err != nil {

					return err
				}
				return nil

			} else {

				position += dao.rehash()
			}

		} else {

			return errors.New("object does not exist to update")
		}
	}
}

func (dao *HashableDao[T]) Delete(object T) error {

	daoIO, err := dao.io()

	defer func() {

		if innerErr := daoIO.Close(); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()
	position := dao.hash(object.Id())

	for {

		occupied, err := daoIO.ReadBoolAt(position)

		if err != nil {

			return err
		}
		if occupied {

			read, err := daoIO.Read()

			if err != nil {

				return err
			}
			if object.Id() == (*read).Id() {

				if err := daoIO.WriteBoolAt(false, position); err != nil {

					return err
				}
				if err := daoIO.Zero(); err != nil {

					return err
				}
				return nil

			} else {

				position += dao.rehash()
			}

		} else {

			return errors.New("object does not exist to delete")
		}
	}
}
