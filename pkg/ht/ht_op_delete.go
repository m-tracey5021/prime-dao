package ht

import "github.com/google/uuid"

func (hashTable *TSFHashTable[T]) Delete(id uuid.UUID) error {

	table, bucket, err := hashTable.tableManager.Locate(id)

	defer hashTable.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return err
	}
	if err := hashTable.fileManager.GoTo(bucket.bucketLocation, table); err != nil {

		return err
	}
	bucketHeader := HashTableBucketHeader{false, true, 0}

	if err := hashTable.bucketHeaderIO.Write(table, bucketHeader); err != nil {

		return err
	}
	if err := hashTable.objectIO.Zero(table); err != nil {

		return err
	}
	return err
}
