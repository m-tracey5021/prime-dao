package ht

func (hashTable *TSFHashTable[T]) Save(object T) error {

	table, emptyBucket, err := hashTable.tableManager.LocateEmpty(object.Id())

	defer hashTable.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return err
	}
	if err := hashTable.fileManager.GoTo(emptyBucket, table); err != nil {

		return err
	}
	bucketHeader := HashTableBucketHeader{true, false, 0}

	if err := hashTable.bucketHeaderIO.Write(table, bucketHeader); err != nil {

		return err
	}
	if err := hashTable.objectIO.Write(table, object); err != nil {

		return err
	}
	return err
}
