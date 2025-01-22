package ht

import (
	"io"
	"os"
	"testing"

	"github.com/m-tracey5021/prime-dao/pkg/data/fm"

	"github.com/stretchr/testify/assert"
)

func TestLocateNoResult(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	hash := tableManager.ComputeHash(id)

	bucketHeader := HashTableBucketHeader{}

	table := new(os.File)

	var innerErr error = nil

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, innerErr)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, nil)

	tableResult, _, err := tableManager.Locate(id)

	// Then
	assert.Nil(t, tableResult)

	assert.Equal(t, ObjectDoesNotExist, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestLocateNoResultWithEOF(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	hash := tableManager.ComputeHash(id)

	bucketHeader := HashTableBucketHeader{}

	table := new(os.File)

	var innerErr error = nil

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, innerErr)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, io.EOF)

	tableResult, _, err := tableManager.Locate(id)

	// Then
	assert.Nil(t, tableResult)

	assert.Equal(t, ObjectDoesNotExist, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestLocateWithResult(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	hash := tableManager.ComputeHash(id)

	bucketHeader := HashTableBucketHeader{occupied: true}

	object := MockHashable{id}

	table := new(os.File)

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, nil)

	mockFileManager.On("CurrentPosition", table).Return(0, nil)

	mockObjectIO.On("Read", table).Return(object, nil)

	tableResult, locationResult, err := tableManager.Locate(id)

	// Then
	assert.Equal(t, table, tableResult)

	assert.Equal(t, object, locationResult.object)

	assert.Nil(t, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestLocateWithIdMismatch(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	mismatchId := uint64(1)

	hash := tableManager.ComputeHash(mismatchId)

	bucketHeader := HashTableBucketHeader{occupied: true}

	object := MockHashable{id}

	table := new(os.File)

	var innerErr error = nil

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, innerErr)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, nil)

	mockFileManager.On("CurrentPosition", table).Return(0, nil)

	mockObjectIO.On("Read", table).Return(object, nil)

	tableResult, _, err := tableManager.Locate(mismatchId)

	// Then
	assert.Nil(t, tableResult)

	assert.Equal(t, ObjectDoesNotExist, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestLocateEmptyNoResultObjectAlreadyExists(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	hash := tableManager.ComputeHash(id)

	bucketHeader := HashTableBucketHeader{occupied: true}

	object := MockHashable{id}

	table := new(os.File)

	var innerErr error = ObjectAlreadyExists

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, nil)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, nil)

	mockObjectIO.On("Read", table).Return(object, nil)

	tableResult, _, err := tableManager.LocateEmpty(id)

	// Then
	assert.Nil(t, tableResult)

	assert.Equal(t, innerErr, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestLocateEmptyNoResultBucketOccupied(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	mismatchId := uint64(1)

	hash := tableManager.ComputeHash(mismatchId)

	bucketHeader := HashTableBucketHeader{occupied: true}

	object := MockHashable{id}

	table := new(os.File)

	var innerErr error = nil

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, innerErr)

	mockFileManager.On("CloseAndUnlock", table, &innerErr).Return(nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, nil)

	mockObjectIO.On("Read", table).Return(object, nil)

	tableResult, _, err := tableManager.LocateEmpty(mismatchId)

	// Then
	assert.Nil(t, tableResult)

	assert.Equal(t, BucketOccupied, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}

func TestLocateEmptyWithResult(t *testing.T) {

	// Given
	mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager := setupTableManagerMockDependencies()

	id := uint64(0)

	hash := tableManager.ComputeHash(id)

	bucketHeader := HashTableBucketHeader{}

	table := new(os.File)

	// When
	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, []uint64{uint64(hash.tableGroup), uint64(hash.tableNumber)}).Return(table, nil)

	mockFileManager.On("GoTo", hash.position, table).Return(nil)

	mockBucketHeaderIO.On("Read", table).Return(bucketHeader, io.EOF)

	tableResult, position, err := tableManager.LocateEmpty(id)

	// Then
	assert.Equal(t, tableResult, table)

	assert.Equal(t, position, hash.position)

	assert.Nil(t, err)

	mockFileManager.AssertExpectations(t)

	mockBucketHeaderIO.AssertExpectations(t)

	mockObjectIO.AssertExpectations(t)
}
