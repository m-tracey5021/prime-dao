package daoIO

import "github.com/stretchr/testify/mock"

type MockDaoIOFactory[T any] struct {
	mock.Mock
}

func (mockFactory *MockDaoIOFactory[T]) Create(filePath string) (IDaoIO[T], error) {

	args := mockFactory.Called(filePath)

	return args.Get(0).(IDaoIO[T]), args.Error(1)
}
