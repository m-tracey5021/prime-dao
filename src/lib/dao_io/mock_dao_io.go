package daoIO

import "github.com/stretchr/testify/mock"

type MockDaoIO[T any] struct {
	mock.Mock
}

func (mockDaoIO *MockDaoIO[T]) Write(object T) (int, error) {

	args := mockDaoIO.Called(object)

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) WriteAt(object T, position int) (int, error) {

	args := mockDaoIO.Called(object, position)

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) Read() (*T, error) {

	args := mockDaoIO.Called()

	return args.Get(0).(*T), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) ReadBool() (bool, error) {

	args := mockDaoIO.Called()

	return args.Bool(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) ReadAt(position int) (*T, error) {

	args := mockDaoIO.Called(position)

	return args.Get(0).(*T), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) ReadBoolAt(position int) (bool, error) {

	args := mockDaoIO.Called(position)

	return args.Bool(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) Update(object T) error {

	args := mockDaoIO.Called(object)

	return args.Error(0)
}

func (mockDaoIO *MockDaoIO[T]) UpdateAt(object T, position int) error {

	args := mockDaoIO.Called(object, position)

	return args.Error(0)
}

func (mockDaoIO *MockDaoIO[T]) Delete() (int, error) {

	args := mockDaoIO.Called()

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) DeleteAt(position int) (int, error) {

	args := mockDaoIO.Called(position)

	return args.Int(0), args.Error(1)
}

func (mockDaoIO *MockDaoIO[T]) Close() error {

	args := mockDaoIO.Called()

	return args.Error(0)
}
