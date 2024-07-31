package dao

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockDaoFileContainer struct {
	mock.Mock
}

func (mockContainer *MockDaoFileContainer) ManagingFile() (*os.File, error) {

	args := mockContainer.Called()

	return args.Get(0).(*os.File), args.Error(1)
}

func (mockContainer *MockDaoFileContainer) MainTable() (*os.File, error) {

	args := mockContainer.Called()

	return args.Get(0).(*os.File), args.Error(1)
}

func (mockContainer *MockDaoFileContainer) CollisionTable(fileId uint64) (*os.File, error) {

	args := mockContainer.Called(fileId)

	return args.Get(0).(*os.File), args.Error(1)
}

func (mockContainer *MockDaoFileContainer) Size(file *os.File) (int, error) {

	args := mockContainer.Called(file)

	return args.Int(0), args.Error(1)
}

func (mockContainer *MockDaoFileContainer) GoTo(position int, file *os.File) error {

	args := mockContainer.Called(position, file)

	return args.Error(0)
}

func (mockContainer *MockDaoFileContainer) CurrentPosition(file *os.File) (int, error) {

	args := mockContainer.Called(file)

	return args.Int(0), args.Error(1)
}

func (mockContainer *MockDaoFileContainer) Close(file *os.File) error {

	args := mockContainer.Called(file)

	return args.Error(0)
}
