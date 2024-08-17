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

func TestUpdate(t *testing.T) {

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

// func TestUpdateAtEmptyPosition(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	objectToUpdate := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{false, false, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	// When
// 	err := dao.Update(objectToUpdate)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object does not exist to update", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestUpdateAtEmptyPositionAndEOF(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	objectToUpdate := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, io.EOF)

// 	// When
// 	err := dao.Update(objectToUpdate)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object does not exist to update", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestUpdateAtOccupiedPosition(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	savedObject := MockHashable{id}

// 	objectToUpdate := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(savedObject, nil)

// 	mockObjectIO.On("Write", mockMainTable, objectToUpdate).Return(nil)

// 	// When
// 	err := dao.Update(objectToUpdate)

// 	// Then
// 	assert.Nil(t, err)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestUpdateAtOccupiedPositionWithCollision(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	savedObject := MockHashable{id}

// 	collision := MockHashable{1}

// 	objectToUpdate := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil).Once()

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

// 	mockFileManager.On("CurrentPosition", mockCollisionTable).Return(0, nil)

// 	mockObjectIO.On("Read", mockCollisionTable).Return(savedObject, nil)

// 	mockFileManager.On("GoTo", 0, mockCollisionTable).Return(nil)

// 	mockObjectIO.On("Write", mockCollisionTable, objectToUpdate).Return(nil)

// 	// When
// 	err := dao.Update(objectToUpdate)

// 	// Then
// 	assert.Nil(t, err)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestUpdateAtOccupiedPositionWithCollisionAndEOF(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	collision := MockHashable{1}

// 	objectToUpdate := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil).Once()

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

// 	mockFileManager.On("CurrentPosition", mockCollisionTable).Return(0, nil)

// 	mockObjectIO.On("Read", mockCollisionTable).Return(MockHashable{}, io.EOF)

// 	// When
// 	err := dao.Update(objectToUpdate)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object does not exist to update", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }
