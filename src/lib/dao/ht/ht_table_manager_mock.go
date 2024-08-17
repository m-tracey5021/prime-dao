package ht

import (
	"os"
	"transformer/src/lib/dao/schema"

	"github.com/stretchr/testify/mock"
)

type MockHashTableManager[T schema.FixedSizeIdentifiable] struct {
	mock.Mock
}

func (mockTableManager *MockHashTableManager[T]) ComputeHash(id uint64) Hash {

	args := mockTableManager.Called(id)

	return args.Get(0).(Hash)
}

func (mockTableManager *MockHashTableManager[T]) Locate(id uint64) (*os.File, HashTableBucket[T], error) {

	args := mockTableManager.Called(id)

	return args.Get(0).(*os.File), args.Get(1).(HashTableBucket[T]), args.Error(2)
}

func (mockTableManager *MockHashTableManager[T]) LocateEmpty(id uint64) (*os.File, int, error) {

	args := mockTableManager.Called(id)

	return args.Get(0).(*os.File), args.Int(1), args.Error(2)
}
