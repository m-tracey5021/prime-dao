package ht

func (ht *TSFHashTable[T]) Save(object T) error {

	table, emptyBucket, err := ht.LocateEmpty(object.Id())

	defer ht.fileContainer.Close(table, &err)

	if err != nil {

		return err
	}
	if err := ht.fileContainer.GoTo(emptyBucket, table); err != nil {

		return err
	}
	bucketHeader := HashTableBucketHeader{true, false, 0}

	if err := ht.bucketHeaderIO.Write(table, bucketHeader); err != nil {

		return err
	}
	if err := ht.objectIO.Write(table, object); err != nil {

		return err
	}
	return err
}
