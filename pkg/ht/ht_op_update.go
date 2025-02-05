package ht

func (hashTable *TSFHashTable[T]) Update(object T) error {

	table, bucket, err := hashTable.tableManager.Locate(object.Id())

	defer hashTable.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return err
	}
	if err := hashTable.fileManager.GoTo(bucket.objectLocation, table); err != nil {

		return err
	}
	if err := hashTable.objectIO.Write(table, object); err != nil {

		return err
	}
	return err
}
