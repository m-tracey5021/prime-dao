package ht

import (
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"

	"github.com/stretchr/testify/mock"
)

type MockHashTableManager[T schema.FixedSizeIdentifiable] struct {
	mock.Mock
}

func (mockTableManager *MockHashTableManager[T]) ComputeHash(id uuid.UUID) Hash {

	args := mockTableManager.Called(id)

	return args.Get(0).(Hash)
}

func (mockTableManager *MockHashTableManager[T]) Locate(id uuid.UUID) (*os.File, HashTableBucket[T], error) {

	args := mockTableManager.Called(id)

	return args.Get(0).(*os.File), args.Get(1).(HashTableBucket[T]), args.Error(2)
}

func (mockTableManager *MockHashTableManager[T]) LocateEmpty(id uuid.UUID) (*os.File, int, error) {

	args := mockTableManager.Called(id)

	return args.Get(0).(*os.File), args.Int(1), args.Error(2)
}
