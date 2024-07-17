package dao

import (
	"errors"
	"testing"
	. "transformer/src/lib/component/hashable_component"
	. "transformer/src/lib/dao_io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDaoIO[T any] struct {
	mock.Mock
}

func (mockDaoIO *MockDaoIO[T]) ReadBoolAt() (bool, error) {

	args := mockDaoIO.Called()

	return args.Bool(0), args.Error(1)
}

type MockDaoIOFactory[T any] struct {
	mock.Mock
}

func (mockFactory *MockDaoIOFactory[T]) Create(filePath string) (IDaoIO[T], error) {

	args := mockFactory.Called()

	return args.Get(0).(*DaoIO[T]), args.Error(1)
}

func TestNewId(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	// When
	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	idA, errA := dao.NewId()

	idB, errB := dao.NewId()

	err := errors.Join(errA, errB)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !(assert.Equal(t, 0, idA) && assert.Equal(t, 1, idB)) {

		t.Fail()
	}
}
