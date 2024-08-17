package daoio

import (
	"os"
	"transformer/src/lib/dao/schema"

	"github.com/stretchr/testify/mock"
)

type MockFixedSizeDao[T schema.FixedSize] struct {
	mock.Mock
}

func (mockDaoIO *MockFixedSizeDao[T]) Write(file *os.File, object T) error {

	args := mockDaoIO.Called(file, object)

	return args.Error(0)
}

func (mockDaoIO *MockFixedSizeDao[T]) Read(file *os.File) (T, error) {

	args := mockDaoIO.Called(file)

	return args.Get(0).(T), args.Error(1)
}

func (mockDaoIO *MockFixedSizeDao[T]) Zero(file *os.File) error {

	args := mockDaoIO.Called(file)

	return args.Error(0)
}
