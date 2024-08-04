package dao

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUpdateAtEmptyPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	id := uint64(0)

	objectToUpdate := MockHashable{id}

	bucketHeader := HashTableBucketHeader{false, false, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	// When
	err := dao.Update(objectToUpdate)

	// Then
	assert.NotNil(t, err)

	assert.Equal(t, "object does not exist to update", err.Error())

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestUpdateAtEmptyPositionAndEOF(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	id := uint64(0)

	objectToUpdate := MockHashable{id}

	bucketHeader := HashTableBucketHeader{}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, io.EOF)

	// When
	err := dao.Update(objectToUpdate)

	// Then
	assert.NotNil(t, err)

	assert.Equal(t, "object does not exist to update", err.Error())

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestUpdateAtOccupiedPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	id := uint64(0)

	savedObject := MockHashable{id}

	objectToUpdate := MockHashable{id}

	bucketHeader := HashTableBucketHeader{true, false, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(savedObject, nil)

	mockObjectIO.On("Write", mockMainTable, objectToUpdate).Return(nil)

	// When
	err := dao.Update(objectToUpdate)

	// Then
	assert.Nil(t, err)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestUpdateAtOccupiedPositionWithCollision(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	id := uint64(0)

	savedObject := MockHashable{id}

	collision := MockHashable{1}

	objectToUpdate := MockHashable{id}

	bucketHeader := HashTableBucketHeader{true, false, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil).Once()

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

	mockFileContainer.On("CurrentPosition", mockCollisionTable).Return(0, nil)

	mockObjectIO.On("Read", mockCollisionTable).Return(savedObject, nil)

	mockFileContainer.On("GoTo", 0, mockCollisionTable).Return(nil)

	mockObjectIO.On("Write", mockCollisionTable, objectToUpdate).Return(nil)

	// When
	err := dao.Update(objectToUpdate)

	// Then
	assert.Nil(t, err)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockFileContainer.AssertCalled(t, "File", HashTableCollisionTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestUpdateAtOccupiedPositionWithCollisionAndEOF(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	id := uint64(0)

	collision := MockHashable{1}

	objectToUpdate := MockHashable{id}

	bucketHeader := HashTableBucketHeader{true, false, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil).Once()

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

	mockFileContainer.On("CurrentPosition", mockCollisionTable).Return(0, nil)

	mockObjectIO.On("Read", mockCollisionTable).Return(MockHashable{}, io.EOF)

	// When
	err := dao.Update(objectToUpdate)

	// Then
	assert.NotNil(t, err)

	assert.Equal(t, "object does not exist to update", err.Error())

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockFileContainer.AssertCalled(t, "File", HashTableCollisionTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}
