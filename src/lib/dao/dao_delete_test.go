package dao

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeleteAtEmptyPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	idToDeleteToDelete := uint64(0)

	mockBucketHeader := DaoBucketHeader{false, false, 0}

	hash := dao.hash(idToDeleteToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	// When
	err := dao.Delete(idToDeleteToDelete)

	// Then
	assert.NotNil(t, err)

	assert.Equal(t, "object does not exist to delete", err.Error())

	mockFileContainer.AssertCalled(t, "MainTable")

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestDeleteAtEmptyPositionAndEOF(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	idToDelete := uint64(0)

	mockBucketHeader := DaoBucketHeader{}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, io.EOF)

	// When
	err := dao.Delete(idToDelete)

	// Then
	assert.NotNil(t, err)

	assert.Equal(t, "object does not exist to delete", err.Error())

	mockFileContainer.AssertCalled(t, "MainTable")

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestDeleteAtOccupiedPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	idToDelete := uint64(0)

	savedObject := MockHashable{idToDelete}

	mockBucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(savedObject, nil)

	mockObjectIO.On("Zero", mockMainTable).Return(nil)

	// When
	err := dao.Delete(idToDelete)

	// Then
	assert.Nil(t, err)

	mockFileContainer.AssertCalled(t, "MainTable")

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestDeleteAtOccupiedPositionWithCollision(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	idToDelete := uint64(0)

	savedObject := MockHashable{idToDelete}

	collision := MockHashable{1}

	mockBucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil).Once()

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

	mockFileContainer.On("CurrentPosition", mockCollisionTable).Return(0, nil)

	mockObjectIO.On("Read", mockCollisionTable).Return(savedObject, nil)

	mockFileContainer.On("GoTo", 0, mockCollisionTable).Return(nil)

	mockObjectIO.On("Zero", mockCollisionTable).Return(nil)

	// When
	err := dao.Delete(idToDelete)

	// Then
	assert.Nil(t, err)

	mockFileContainer.AssertCalled(t, "MainTable")

	mockFileContainer.AssertCalled(t, "CollisionTable", mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestDeleteAtOccupiedPositionWithCollisionAndEOF(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	idToDelete := uint64(0)

	collision := MockHashable{1}

	mockBucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil).Once()

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

	mockFileContainer.On("CurrentPosition", mockCollisionTable).Return(0, nil)

	mockObjectIO.On("Read", mockCollisionTable).Return(MockHashable{}, io.EOF)

	// When
	err := dao.Delete(idToDelete)

	// Then
	assert.NotNil(t, err)

	assert.Equal(t, "object does not exist to delete", err.Error())

	mockFileContainer.AssertCalled(t, "MainTable")

	mockFileContainer.AssertCalled(t, "CollisionTable", mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}
