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

	mockBucketHeader := DaoBucketHeader{}

	hash := dao.hash(id)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, io.EOF)

	// When
	read, err := dao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Nil(t, read)

	mockFileContainer.AssertCalled(t, "MainTable")

	mockBucketHeaderIO.AssertExpectations(t)
}

func TestGetAtOccupiedPosition(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, _ := setupMockFiles(mockFileContainer)

	mockId := uint64(0)

	mockObjectA := MockHashable{mockId}

	mockBucketHeader := DaoBucketHeader{true, false, 0}

	hash := dao.hash(mockId)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(mockObjectA, nil)

	// When
	read, err := dao.Get(mockId)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Equal(t, &mockObjectA, read)

	mockFileContainer.AssertCalled(t, "MainTable")

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestGetAtOccupiedPositionWithCollision(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	mockId := uint64(0)

	mockObjectA := MockHashable{1}

	mockObjectB := MockHashable{mockId}

	mockBucketHeader := DaoBucketHeader{true, true, 0}

	hash := dao.hash(mockId)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(mockObjectA, nil).Once()

	mockObjectIO.On("Read", mockCollisionTable).Return(mockObjectB, nil)

	// When
	read, err := dao.Get(mockId)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Equal(t, &mockObjectB, read)

	mockFileContainer.AssertCalled(t, "MainTable")

	mockFileContainer.AssertCalled(t, "CollisionTable", mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestGetAtOccupiedPositionWithCollisionNoResult(t *testing.T) {

	// Given
	mockFileContainer, _, mockBucketHeaderIO, mockObjectIO, dao := setupMockDao()

	_, mockMainTable, mockCollisionTable := setupMockFiles(mockFileContainer)

	mockId := uint64(0)

	mockObjectA := MockHashable{1}

	mockObjectB := MockHashable{mockId}

	mockBucketHeader := DaoBucketHeader{true, true, 0}

	hash := dao.hash(mockId)

	mockFileContainer.On("GoTo", hash, mockMainTable).Return(nil)

	mockBucketHeaderIO.On("Read", mockMainTable).Return(mockBucketHeader, nil)

	mockObjectIO.On("Read", mockMainTable).Return(mockObjectA, nil).Once()

	mockObjectIO.On("Read", mockCollisionTable).Return(mockObjectB, io.EOF)

	// When
	read, err := dao.Get(mockId)

	if err != nil {

		t.Fatalf("%v", err)
	}

	// Then
	assert.Nil(t, read)

	mockFileContainer.AssertCalled(t, "MainTable")

	mockFileContainer.AssertCalled(t, "CollisionTable", mock.AnythingOfType("uint64"))

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}
