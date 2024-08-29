package fm

import (
	"os"

	"github.com/stretchr/testify/mock"
)

type MockFileManager struct {
	mock.Mock
}

func (mockFileManager *MockFileManager) Path() string {

	args := mockFileManager.Called()

	return args.String(0)
}

func (mockFileManager *MockFileManager) Descriptor() string {

	args := mockFileManager.Called()

	return args.String(0)
}

func (mockFileManager *MockFileManager) OpenAndLock(fileAlias FileAlias, idChain ...uint64) (*os.File, error) {

	args := mockFileManager.Called(fileAlias, idChain)

	return args.Get(0).(*os.File), args.Error(1)
}

func (mockFileManager *MockFileManager) Size(file *os.File) (int, error) {

	args := mockFileManager.Called(file)

	return args.Int(0), args.Error(1)
}

func (mockFileManager *MockFileManager) GoTo(position int, file *os.File) error {

	args := mockFileManager.Called(position, file)

	return args.Error(0)
}

func (mockFileManager *MockFileManager) CurrentPosition(file *os.File) (int, error) {

	args := mockFileManager.Called(file)

	return args.Int(0), args.Error(1)
}

func (mockFileManager *MockFileManager) Remove(file *os.File) error {

	args := mockFileManager.Called(file)

	return args.Error(0)
}

func (mockFileManager *MockFileManager) CloseAndUnlock(file *os.File, err *error) error {

	args := mockFileManager.Called(file, err)

	return args.Error(0)
}
