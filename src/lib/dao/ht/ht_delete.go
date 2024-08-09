package ht

func (ht *TSFHashTable[T]) Delete(id uint64) error {

	table, bucket, err := ht.Locate(id)

	defer ht.fileContainer.Close(table, &err)

	if err != nil {

		return err
	}
	if err := ht.fileContainer.GoTo(bucket.bucketLocation, table); err != nil {

		return err
	}
	bucketHeader := HashTableBucketHeader{false, true, 0}

	if err := ht.bucketHeaderIO.Write(table, bucketHeader); err != nil {

		return err
	}
	if err := ht.objectIO.Zero(table); err != nil {

		return err
	}
	return err
}
