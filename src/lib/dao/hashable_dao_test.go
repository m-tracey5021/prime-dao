package dao

import (
	"testing"
	. "transformer/src/lib/component/hashable_component"
	. "transformer/src/lib/dao_io"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNewIdAtEmptyPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", 25).Return(false, nil)

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	// When
	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	id, err := dao.NewId()

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, 0, int(id)) {

		t.Fail()
	}
}

func TestNewIdAtOccupiedPosition(t *testing.T) { // TODO

	// Given
	mockFilePath := "test"

	mockIndex := NewIndex(0, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", 25).Return(true, nil)

	mockDaoIO.On("ReadBoolAt", 27).Return(false, nil)

	mockDaoIO.On("ReadAt", 25).Return(&mockIndex, nil)

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	// When
	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	id, err := dao.NewId()

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, 1, int(id)) {

		t.Fail()
	}
}

func TestGetAtEmptyPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", 25).Return(false, nil)

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	// When
	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	read, err := dao.Get(0)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Nil(t, read) {

		t.Fail()
	}
}

func TestGetAtOccupiedPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockIndex := NewIndex(0, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", 25).Return(true, nil)

	mockDaoIO.On("ReadAt", 25).Return(&mockIndex, nil)

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	// When
	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	read, err := dao.Get(0)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, &mockIndex, read) {

		t.Fail()
	}
}

func TestGetAtOccupiedPositionWithCollision(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockIndexA := NewIndex(1, 0, 0)
	mockIndexB := NewIndex(2, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", mock.Anything).Return(true, nil)

	mockDaoIO.On("ReadAt", 29).Return(&mockIndexA, nil)
	mockDaoIO.On("ReadAt", 34).Return(&mockIndexB, nil)

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	// When
	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	read, err := dao.Get(2)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, &mockIndexB, read) {

		t.Fail()
	}
}
