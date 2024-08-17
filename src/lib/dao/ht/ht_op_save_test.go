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

func TestSave(t *testing.T) {

	// Given
	mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	id := uint64(0)

	objectToSave := MockHashable{id}

	table := new(os.File)

	emptyBucket := 0

	var innerErr error = nil

	// When
	mockTableManager.On("LocateEmpty", id).Return(table, emptyBucket, nil)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", emptyBucket, table).Return(nil)

	mockBucketHeaderIO.On("Write", table, mock.AnythingOfType("HashTableBucketHeader")).Return(nil)

	mockObjectIO.On("Write", table, objectToSave).Return(nil)

	err := dao.Save(objectToSave)

	// Then
	assert.Nil(t, err)

	mockFileManager.AssertExpectations(t)

	mockTableManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

// func TestSaveAtOccupiedPositionFirstCollision(t *testing.T) {

// 	// Given
// 	mockFileManager, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	mockManagingFile, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	objectToSave := MockHashable{0}

// 	collision := MockHashable{1}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(objectToSave.Id())

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collision, nil).Once()

// 	mockObjectIO.On("Read", mockCollisionTable).Return(MockHashable{}, io.EOF)

// 	mockObjectIO.On("Write", mockCollisionTable, objectToSave).Return(nil)

// 	mockBucketHeaderIO.On("Write", mockMainTable, mock.AnythingOfType("DaoBucketHeader")).Return(nil)

// 	mockCacheIO.On("WriteSizePrefixed", mockManagingFile, mock.AnythingOfType("DaoIdentifierCache")).Return(0, nil)

// 	// When
// 	err := dao.Save(objectToSave)

// 	// Then
// 	assert.Nil(t, err)

// 	mockFileManager.AssertExpectations(t)

// 	mockCacheIO.AssertExpectations(t)

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestSaveAtOccupiedPositionFirstCollisionWithExistingId(t *testing.T) {

// 	// Given
// 	mockFileManager, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	objectToSave := MockHashable{0}

// 	collision := MockHashable{0}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(objectToSave.Id())

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collision, nil)

// 	// When
// 	err := dao.Save(objectToSave)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object already exists, cannot save new", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockCacheIO.AssertExpectations(t)

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestSaveAtOccupiedPositionNonFirstCollision(t *testing.T) {

// 	// Given
// 	mockFileManager, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	mockManagingFile, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	objectToSave := MockHashable{0}

// 	collisionA := MockHashable{1}

// 	collisionB := MockHashable{2}

// 	bucketHeader := HashTableBucketHeader{true, true, 0}

// 	hash := dao.hash(objectToSave.Id())

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collisionA, nil).Once()

// 	mockObjectIO.On("Read", mockCollisionTable).Return(collisionB, nil).Once()

// 	mockObjectIO.On("Read", mockCollisionTable).Return(MockHashable{}, io.EOF)

// 	mockObjectIO.On("Write", mockCollisionTable, objectToSave).Return(nil)

// 	mockCacheIO.On("WriteSizePrefixed", mockManagingFile, mock.AnythingOfType("DaoIdentifierCache")).Return(0, nil)

// 	// When
// 	err := dao.Save(objectToSave)

// 	// Then
// 	assert.Nil(t, err)

// 	mockFileManager.AssertExpectations(t)

// 	mockCacheIO.AssertExpectations(t)

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestSaveAtOccupiedPositionNonFirstCollisionWithExistingId(t *testing.T) {

// 	// Given
// 	mockFileManager, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	objectToSave := MockHashable{0}

// 	collisionA := MockHashable{1}

// 	collisionB := MockHashable{0}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(objectToSave.Id())

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(collisionA, nil).Once()

// 	mockObjectIO.On("Read", mockCollisionTable).Return(collisionB, nil).Once()

// 	// When
// 	err := dao.Save(objectToSave)

// 	// Then
// 	assert.NotNil(t, err)

// 	assert.Equal(t, "object already exists, cannot save new", err.Error())

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableCollisionTable, mock.AnythingOfType("uint64"))

// 	mockCacheIO.AssertExpectations(t)

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }
