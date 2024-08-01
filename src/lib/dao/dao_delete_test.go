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

	idToDelete := uint64(0)

	bucketHeader := DaoBucketHeader{false, false, 0}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	// When
	err := dao.Delete(idToDelete)

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

	bucketHeader := DaoBucketHeader{}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, io.EOF)

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
	mockFileContainer, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	mockManagingFile, mockMainTable, _ := setupMockFiles(mockFileContainer)

	idToDelete := uint64(0)

	savedObject := MockHashable{idToDelete}

	bucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(idToDelete)

	dao.identifierCache.ObjectIds = append(dao.identifierCache.ObjectIds, idToDelete)

	cacheAfterDeletion := DaoIdentifierCache{make([]uint64, 0), make([]uint64, 0)}

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(savedObject, nil)

	mockBucketHeaderIO.On("Write", mockMainTable, mock.AnythingOfType("DaoBucketHeader")).Return(nil)

	mockObjectIO.On("Zero", mockMainTable).Return(nil)

	mockCacheIO.On("WriteSizePrefixed", mockManagingFile, cacheAfterDeletion).Return(0, nil)

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
	mockFileContainer, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	mockManagingFile, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	idToDelete := uint64(0)

	savedObject := MockHashable{idToDelete}

	collision := MockHashable{1}

	bucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(idToDelete)

	dao.identifierCache.ObjectIds = append(dao.identifierCache.ObjectIds, idToDelete)

	cacheAfterDeletion := DaoIdentifierCache{make([]uint64, 0), make([]uint64, 0)}

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil).Once()

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

	mockFileContainer.On("CurrentPosition", mockCollisionTable).Return(0, nil)

	mockObjectIO.On("Read", mockCollisionTable).Return(savedObject, nil)

	mockFileContainer.On("GoTo", 0, mockCollisionTable).Return(nil)

	mockObjectIO.On("Delete", mockCollisionTable).Return(nil)

	mockCacheIO.On("WriteSizePrefixed", mockManagingFile, cacheAfterDeletion).Return(0, nil)

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

	bucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(idToDelete)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil).Once()

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

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
