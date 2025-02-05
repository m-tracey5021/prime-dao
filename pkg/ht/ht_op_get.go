package ht

import "github.com/google/uuid"

func (hashTable *TSFHashTable[T]) Get(id uuid.UUID) (*T, error) {

	table, bucket, err := hashTable.tableManager.Locate(id)

	defer hashTable.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return nil, err
	}
	return &bucket.object, err
}
