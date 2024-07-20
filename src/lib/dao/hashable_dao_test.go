package dao

import (
	"io"
	"testing"
	. "transformer/src/lib/component/hashable_component"
	. "transformer/src/lib/dao_io"

	"github.com/stretchr/testify/assert"
)

func TestNewIdAtEmptyPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockId := 0

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(uint64(mockId))

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(false, nil)

	// When
	id, err := dao.NewId()

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, 0, int(id)) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestNewIdAtOccupiedPosition(t *testing.T) { // TODO

	// Given
	mockFilePath := "test"

	mockIdA := 0

	mockIdB := 1

	mockIndex := Index{}

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(uint64(mockIdA))

	secondHash := dao.hash(uint64(mockIdB))

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(true, nil)

	mockDaoIO.On("ReadBoolAt", secondHash).Return(false, nil)

	mockDaoIO.On("ReadAt", firstHash).Return(&mockIndex, nil)

	// When
	id, err := dao.NewId()

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, 1, int(id)) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestGetAtEmptyPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockId := 0

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	hash := dao.hash(uint64(mockId))

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", hash).Return(false, nil)

	// When
	read, err := dao.Get(0)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Nil(t, read) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestGetAtOccupiedPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockId := 2

	mockIndexA := NewIndex(1, 0, 0)

	mockIndexB := NewIndex(uint64(mockId), 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(uint64(mockId))

	secondHash := firstHash + dao.rehash()

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(true, nil)

	mockDaoIO.On("ReadBoolAt", secondHash).Return(true, nil)

	mockDaoIO.On("ReadAt", firstHash).Return(&mockIndexA, nil)

	mockDaoIO.On("ReadAt", secondHash).Return(&mockIndexB, nil)

	// When
	read, err := dao.Get(uint64(mockId))

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !assert.Equal(t, &mockIndexB, read) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestGetAll(t *testing.T) {

	// Given
	mockFilePath := "test"

	mockIndexA := NewIndex(0, 0, 0)

	mockIndexB := NewIndex(1, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstPosition := 0

	secondPosition := firstPosition + dao.bucketSize

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstPosition).Return(false, nil)

	mockDaoIO.On("ReadBoolAt", secondPosition).Return(true, nil)

	mockDaoIO.On("Read").Return(&mockIndexA, nil).Once()

	mockDaoIO.On("Read").Return(&mockIndexB, nil).Once()

	mockDaoIO.On("Read").Return(&mockIndexB, io.EOF)

	// When
	results, err := dao.GetAll()

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	if !(assert.Equal(t, 2, len(results)) &&
		assert.Equal(t, mockIndexA, results[0]) &&
		assert.Equal(t, mockIndexB, results[1])) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestSaveAtOccupiedPosition(t *testing.T) {

	// Given
	indexToSave := NewIndex(1, 0, 0)

	mockFilePath := "test"

	mockIndex := NewIndex(0, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(indexToSave.Id())

	secondHash := firstHash + dao.rehash()

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("Read").Return(&mockIndex, nil)

	mockDaoIO.On("Write", indexToSave).Return(0, nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(true, nil)

	mockDaoIO.On("ReadBoolAt", secondHash).Return(false, nil)

	mockDaoIO.On("WriteBoolAt", true, secondHash).Return(nil)

	// When
	err := dao.Save(indexToSave)

	// Then
	if !assert.Nil(t, err) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestSaveWhereObjectAlreadyExists(t *testing.T) {

	// Given
	indexToSave := NewIndex(0, 0, 0)

	mockFilePath := "test"

	mockIndex := NewIndex(0, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	hash := dao.hash(indexToSave.Id())

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("Read").Return(&mockIndex, nil)

	mockDaoIO.On("ReadBoolAt", hash).Return(true, nil)

	// When
	err := dao.Save(indexToSave)

	// Then
	if !(assert.NotNil(t, err) && assert.Equal(t, "object already exists, cannot save new", err.Error())) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestUpdateAtOccupiedPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	indexToUpdate := NewIndex(1, 0, 0)

	mockIndexA := NewIndex(0, 0, 0)
	mockIndexB := NewIndex(1, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(indexToUpdate.Id())

	secondHash := firstHash + dao.rehash()

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(true, nil)

	mockDaoIO.On("ReadBoolAt", secondHash).Return(true, nil)

	mockDaoIO.On("Read").Return(&mockIndexA, nil).Once()

	mockDaoIO.On("Read").Return(&mockIndexB, nil) // this should be the second call

	mockDaoIO.On("UpdateAt", indexToUpdate, secondHash+1).Return(nil)

	// When
	err := dao.Update(indexToUpdate)

	// Then
	if !assert.Nil(t, err) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestUpdateWhereObjectDoesNotExist(t *testing.T) {

	// Given
	mockFilePath := "test"

	indexToUpdate := NewIndex(1, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(indexToUpdate.Id())

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(false, nil)

	// When
	err := dao.Update(indexToUpdate)

	// Then
	if !(assert.NotNil(t, err) && assert.Equal(t, "object does not exist to update", err.Error())) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestDeleteAtOccupiedPosition(t *testing.T) {

	// Given
	mockFilePath := "test"

	indexToDelete := NewIndex(1, 0, 0)

	mockIndexA := NewIndex(0, 0, 0)
	mockIndexB := NewIndex(1, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(indexToDelete.Id())

	secondHash := firstHash + dao.rehash()

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(true, nil)

	mockDaoIO.On("ReadBoolAt", secondHash).Return(true, nil)

	mockDaoIO.On("Read").Return(&mockIndexA, nil).Once()

	mockDaoIO.On("Read").Return(&mockIndexB, nil) // this should be the second call

	mockDaoIO.On("WriteBoolAt", false, secondHash).Return(nil)

	mockDaoIO.On("Zero").Return(nil)

	// When
	err := dao.Delete(indexToDelete)

	// Then
	if !assert.Nil(t, err) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}

func TestDeleteWhereObjectDoesNotExist(t *testing.T) {

	// Given
	mockFilePath := "test"

	indexToDelete := NewIndex(1, 0, 0)

	mockDaoIO := new(MockDaoIO[Index])

	mockDaoIOFactory := new(MockDaoIOFactory[Index])

	mockDaoIOFactory.On("Create", mockFilePath).Return(mockDaoIO, nil)

	dao := NewHashableDao[Index](mockDaoIOFactory, mockFilePath, 10)

	firstHash := dao.hash(indexToDelete.Id())

	mockDaoIO.On("Close").Return(nil)

	mockDaoIO.On("ReadBoolAt", firstHash).Return(false, nil)

	// When
	err := dao.Delete(indexToDelete)

	// Then
	if !(assert.NotNil(t, err) && assert.Equal(t, "object does not exist to delete", err.Error())) {

		t.Fail()
	}
	mockDaoIO.AssertExpectations(t)

	mockDaoIOFactory.AssertExpectations(t)
}
