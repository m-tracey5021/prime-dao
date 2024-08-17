package ht

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// import (
// 	"io"
// 	"testing"
// 	"transformer/src/lib/dao/fm"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

func TestDelete(t *testing.T) {

	// Given
	mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	id := uint64(0)

	table := new(os.File)

	location := 0

	bucket := HashTableBucket[MockHashable]{bucketLocation: location}

	var innerErr error = nil

	// When
	mockTableManager.On("Locate", id).Return(table, bucket, nil)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", bucket.bucketLocation, table).Return(nil)

	mockBucketHeaderIO.On("Write", table, mock.AnythingOfType("HashTableBucketHeader")).Return(nil)

	mockObjectIO.On("Zero", table).Return(nil)

	err := dao.Delete(id)

	// Then
	assert.Nil(t, err)

	mockFileManager.AssertExpectations(t)

	mockTableManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

// func TestDeleteAtEmptyPosition(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	idToDelete := uint64(0)

// 	bucketHeader := HashTableBucketHeader{false, false, 0}

// 	hash := dao.hash(idToDelete)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	// When
// 	err := dao.Delete(idToDelete)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object does not exist to delete", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestDeleteAtEmptyPositionAndEOF(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	idToDelete := uint64(0)

// 	bucketHeader := HashTableBucketHeader{}

// 	hash := dao.hash(idToDelete)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, io.EOF)

// 	// When
// 	err := dao.Delete(idToDelete)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object does not exist to delete", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestDeleteAtOccupiedPosition(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	idToDelete := uint64(0)

// 	savedObject := MockHashable{idToDelete}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(idToDelete)

// 	// dao.identifierCache.ObjectIds = append(dao.identifierCache.ObjectIds, idToDelete)

// 	// cacheAfterDeletion := HashTableIdentifierCache{make([]uint64, 0), make([]uint64, 0)}

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(savedObject, nil)

// 	mockBucketHeaderIO.On("Write", mockMainTable, mock.AnythingOfType("DaoBucketHeader")).Return(nil)

// 	mockObjectIO.On("Zero", mockMainTable).Return(nil)

// 	// mockCacheIO.On("WriteSizePrefixed", mockManagingFile, cacheAfterDeletion).Return(0, nil)

// 	// When
// 	err := dao.Delete(idToDelete)

// 	// Then
// 	assert.Nil(t, err)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestDeleteAtOccupiedPositionWithCollision(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	idToDelete := uint64(0)

// 	savedObject := MockHashable{idToDelete}

// 	collision := MockHashable{1}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(idToDelete)

// 	// dao.identifierCache.ObjectIds = append(dao.identifierCache.ObjectIds, idToDelete)

// 	// cacheAfterDeletion := HashTableIdentifierCache{make([]uint64, 0), make([]uint64, 0)}

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil).Once()

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

// 	mockFileManager.On("CurrentPosition", mockCollisionTable).Return(0, nil)

// 	mockObjectIO.On("Read", mockCollisionTable).Return(savedObject, nil)

// 	mockFileManager.On("GoTo", 0, mockCollisionTable).Return(nil)

// 	mockObjectIO.On("Delete", mockCollisionTable).Return(nil)

// 	// mockCacheIO.On("WriteSizePrefixed", mockManagingFile, cacheAfterDeletion).Return(0, nil)

// 	// When
// 	err := dao.Delete(idToDelete)

// 	// Then
// 	assert.Nil(t, err)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableCollisionTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestDeleteAtOccupiedPositionWithCollisionAndEOF(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	idToDelete := uint64(0)

// 	collision := MockHashable{1}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(idToDelete)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil).Once()

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

// 	mockFileManager.On("CurrentPosition", mockCollisionTable).Return(0, nil)

// 	mockObjectIO.On("Read", mockCollisionTable).Return(MockHashable{}, io.EOF)

// 	// When
// 	err := dao.Delete(idToDelete)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object does not exist to delete", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableCollisionTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }
