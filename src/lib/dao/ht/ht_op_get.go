package ht

func (hashTable *TSFHashTable[T]) Get(id uint64) (*T, error) {

	table, bucket, err := hashTable.tableManager.Locate(id)

	defer hashTable.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return nil, err
	}
	return &bucket.object, err
}
