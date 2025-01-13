package dataio

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockDataIO[T any] struct {
	mock.Mock
}

func (mockDataIO *MockDataIO[T]) WriteSizePrefixed(file *os.File, object T) (int, error) {

	args := mockDataIO.Called(file, object)

	return args.Int(0), args.Error(1)
}

func (mockDataIO *MockDataIO[T]) ReadSizePrefixed(file *os.File) (T, error) {

	args := mockDataIO.Called(file)

	return args.Get(0).(T), args.Error(1)
}

func (mockDataIO *MockDataIO[T]) Update(file *os.File, object T) (int, error) {

	args := mockDataIO.Called(file, object)

	return args.Int(0), args.Error(1)
}

func (mockDataIO *MockDataIO[T]) Delete(file *os.File) (int, error) {

	args := mockDataIO.Called(file)

	return args.Int(0), args.Error(1)
}

func (mockDataIO *MockDataIO[T]) Zero(file *os.File) error {

	args := mockDataIO.Called(file)

	return args.Error(0)
}
