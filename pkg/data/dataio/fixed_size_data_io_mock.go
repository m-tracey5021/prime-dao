package dataio

import (
	"os"
	"prime-dao/pkg/data/schema"

	"github.com/stretchr/testify/mock"
)

type MockFixedSizeDataIO[T schema.FixedSize] struct {
	mock.Mock
}

func (mockDataIO *MockFixedSizeDataIO[T]) Write(file *os.File, object T) error {

	args := mockDataIO.Called(file, object)

	return args.Error(0)
}

func (mockDataIO *MockFixedSizeDataIO[T]) Read(file *os.File) (T, error) {

	args := mockDataIO.Called(file)

	return args.Get(0).(T), args.Error(1)
}

func (mockDataIO *MockFixedSizeDataIO[T]) Zero(file *os.File) error {

	args := mockDataIO.Called(file)

	return args.Error(0)
}
