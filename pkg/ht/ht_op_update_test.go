package ht

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {

	// Given
	mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, hashTable := setupHashTableMockDependencies()

	id := uuid.New()

	table := new(os.File)

	location := 0

	bucket := HashTableBucket[MockHashable]{bucketLocation: location}

	object := MockHashable{id, 0}

	var innerErr error = nil

	// When
	mockTableManager.On("Locate", id).Return(table, bucket, nil)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", bucket.bucketLocation, table).Return(nil)

	mockObjectIO.On("Write", table, object).Return(nil)

	err := hashTable.Update(object)

	// Then
	assert.Nil(t, err)

	mockFileManager.AssertExpectations(t)

	mockTableManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}
