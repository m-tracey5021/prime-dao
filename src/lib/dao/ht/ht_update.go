package ht

func (ht *TSFHashTable[T]) Update(object T) error {

	table, bucket, err := ht.Locate(object.Id())

	defer ht.fileContainer.Close(table, &err)

	if err != nil {

		return err
	}
	if err := ht.fileContainer.GoTo(bucket.objectLocation, table); err != nil {

		return err
	}
	if err := ht.objectIO.Write(table, object); err != nil {

		return err
	}
	return err
}
