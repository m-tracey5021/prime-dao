package ht

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {

	// Given
	mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, dao := setupHashTableMockDependencies()

	id := uuid.New()

	table := new(os.File)

	retrievedObject := MockHashable{id, 0}

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
