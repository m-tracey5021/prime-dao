package daoIO

import "os"

/*
This interface exists so that it can be passed to a dao
and can manage creation of a DaoIO, which can then be mocked
*/

type IDaoIOFactory[T any] interface {
	Create(filePath string) (IDaoIO[T], error)
}

type DaoIOFactory[T any] struct{}

func (factory *DaoIOFactory[T]) Create(filePath string) (IDaoIO[T], error) {

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {

		return nil, err
	}
	return &DaoIO[T]{file}, nil
}
