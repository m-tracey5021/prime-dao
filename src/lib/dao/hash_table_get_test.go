package dao

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetAtEmptyPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, _, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	id := uint64(0)

	bucketHeader := HashTableBucketHeader{}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, io.EOF)

	// When
	read, err := dao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Nil(t, read)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)
}

func TestGetAtEmptyPositionPreviouslyDeleted(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	id := uint64(0)

	objectToRetrieve := MockHashable{id}

	bucketHeader := HashTableBucketHeader{false, true, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockCollisionTable).Return(objectToRetrieve, nil)

	// When
	read, err := dao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Equal(t, &objectToRetrieve, read)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)
}

func TestGetAtOccupiedPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	id := uint64(0)

	objectToRetrieveA := MockHashable{id}

	bucketHeader := HashTableBucketHeader{true, false, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(objectToRetrieveA, nil)

	// When
	read, err := dao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Equal(t, &objectToRetrieveA, read)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestGetAtOccupiedPositionWithCollision(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	id := uint64(0)

	objectToRetrieveA := MockHashable{1}

	objectToRetrieveB := MockHashable{id}

	bucketHeader := HashTableBucketHeader{true, true, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(objectToRetrieveA, nil).Once()

	mockObjectIO.On("Read", mockCollisionTable).Return(objectToRetrieveB, nil)

	// When
	read, err := dao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Equal(t, &objectToRetrieveB, read)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockFileContainer.AssertCalled(t, "File", HashTableCollisionTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestGetAtOccupiedPositionWithCollisionNoResult(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	id := uint64(0)

	objectToRetrieveA := MockHashable{1}

	objectToRetrieveB := MockHashable{id}

	bucketHeader := HashTableBucketHeader{true, true, 0}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(bucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(objectToRetrieveA, nil).Once()

	mockObjectIO.On("Read", mockCollisionTable).Return(objectToRetrieveB, io.EOF)

	// When
	read, err := dao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Nil(t, read)

	mockFileContainer.AssertCalled(t, "File", HashTableMainTable, mock.AnythingOfType("uint64"))

	mockFileContainer.AssertCalled(t, "File", HashTableCollisionTable, mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}
