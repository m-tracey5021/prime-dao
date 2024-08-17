package ht

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// import (
// 	"io"
// 	"testing"
// 	"transformer/src/lib/dao/fm"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/mock"
// )

func TestGet(t *testing.T) {

	// Given
	mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	id := uint64(0)

	table := new(os.File)

	retrievedObject := MockHashable{id}

	bucket := HashTableBucket[MockHashable]{object: retrievedObject}

	var innerErr error = nil

	// When
	mockTableManager.On("Locate", id).Return(table, bucket, nil)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	object, err := dao.Get(id)

	if err != nil {

		t.Fail()
	}
	// Then
	assert.Equal(t, retrievedObject, *object)

	mockFileManager.AssertExpectations(t)

	mockTableManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

// func TestGetAtEmptyPosition(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, _, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	bucketHeader := HashTableBucketHeader{}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, io.EOF)

// 	// When
// 	read, err := dao.Get(id)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}

// 	// Then
// 	assert.Nil(t, read)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)
// }

// func TestGetAtEmptyPositionPreviouslyDeleted(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	objectToRetrieve := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{false, true, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockCollisionTable).Return(objectToRetrieve, nil)

// 	// When
// 	read, err := dao.Get(id)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}

// 	// Then
// 	assert.Equal(t, &objectToRetrieve, read)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)
// }

// func TestGetAtOccupiedPosition(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, _ := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	objectToRetrieveA := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(objectToRetrieveA, nil)

// 	// When
// 	read, err := dao.Get(id)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}

// 	// Then
// 	assert.Equal(t, &objectToRetrieveA, read)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestGetAtOccupiedPositionWithCollision(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	objectToRetrieveA := MockHashable{1}

// 	objectToRetrieveB := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{true, true, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(objectToRetrieveA, nil).Once()

// 	mockObjectIO.On("Read", mockCollisionTable).Return(objectToRetrieveB, nil)

// 	// When
// 	read, err := dao.Get(id)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}

// 	// Then
// 	assert.Equal(t, &objectToRetrieveB, read)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableCollisionTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }

// func TestGetAtOccupiedPositionWithCollisionNoResult(t *testing.T) {

// 	// Given
// 	mockFileManager, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

// 	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileManager)

// 	id := uint64(0)

// 	objectToRetrieveA := MockHashable{1}

// 	objectToRetrieveB := MockHashable{id}

// 	bucketHeader := HashTableBucketHeader{true, true, 0}

// 	hash := dao.hash(id)

// 	mockFileManager.On("GoTo", hash, mockMainTable).Return(nil)

// 	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

// 	mockObjectIO.On("Read", mockMainTable).Return(objectToRetrieveA, nil).Once()

// 	mockObjectIO.On("Read", mockCollisionTable).Return(objectToRetrieveB, io.EOF)

// 	// When
// 	read, err := dao.Get(id)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}

// 	// Then
// 	assert.Nil(t, read)

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableMainTable, mock.AnythingOfType("uint64"))

// 	mockFileManager.AssertCalled(t, "File", fm.HashTableCollisionTable, mock.AnythingOfType("uint64"))

// 	mockBucketHeaderIO.AssertExpectations(t)

// 	mockObjectIO.AssertExpectations(t)
// }
