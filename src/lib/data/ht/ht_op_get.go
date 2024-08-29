package ht

func (hashTable *TSFHashTable[T]) Get(id uint64) (*T, error) {

	table, bucket, err := hashTable.tableManager.Locate(id)

	defer hashTable.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return nil, err
	}
	return &bucket.object, err
}

func (hashTable *TSFHashTable[T]) GetSome(ids ...uint64) ([]*T, error) {

	results := make([]*T, 0)

	for _, id := range ids {

		object, err := hashTable.Get(id)

		if err != nil {

			return results, err
		}
		results = append(results, object)
	}
	return results, nil
}
