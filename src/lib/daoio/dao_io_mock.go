package daoio

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockDaoIO[T any] struct {
	mock.Mock
}

func (mockDaoIO *MockDaoIO[T]) WriteSizePrefixed(file *os.File, object T) (int, error) {

	args := mockDaoIO.Called(file, object)

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) ReadSizePrefixed(file *os.File) (*T, error) {

	args := mockDaoIO.Called(file)

	return args.Get(0).(*T), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) Update(file *os.File, object T) (int, error) {

	args := mockDaoIO.Called(file, object)

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) Delete(file *os.File) (int, error) {

	args := mockDaoIO.Called(file)

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) Zero(file *os.File) error {

	args := mockDaoIO.Called(file)

	return args.Error(0)
}

func (mockDaoIO *MockDaoIO[T]) Size(file *os.File) (int, error) {

	args := mockDaoIO.Called(file)

	return args.Int(0), args.Error(1)
}
