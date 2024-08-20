package ht

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSave(t *testing.T) {

	// Given
	mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, dao := setupHashTableMockDependencies()

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
