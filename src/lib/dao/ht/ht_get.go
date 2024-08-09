package ht

func (ht *TSFHashTable[T]) Get(id uint64) (*T, error) {

	table, bucket, err := ht.Locate(id)

	defer ht.fileContainer.Close(table, &err)

	if err != nil {

		return nil, err
	}
	return &bucket.object, nil
}
