package ht

func (ht *TSFHashTable[T]) Save(object T) error {

	id := ht.saveQueue.NewRequestId()

	request := HashTableSaveRequest[T]{id, object}

	return ht.saveQueue.Processor().Process(&request).Error()
}

func (ht *TSFHashTable[T]) QueueSave(object T) {

	id := ht.saveQueue.NewRequestId()

	request := HashTableSaveRequest[T]{id, object}

	ht.saveQueue.AddTask(&request)
}

func (ht *TSFHashTable[T]) QueueSaves(objects ...T) {

	for _, object := range objects {

		id := ht.saveQueue.NewRequestId()

		request := HashTableSaveRequest[T]{id, object}

		ht.saveQueue.AddTask(&request)
	}
}

// func (ht *TSFHashTable[T]) Save(object T) error {

// 	table, emptyBucket, err := ht.LocateEmpty(object.Id())

// 	defer ht.fileManager.Close(table, &err)

// 	if err != nil {

// 		return err
// 	}
// 	if err := ht.fileManager.GoTo(emptyBucket, table); err != nil {

// 		return err
// 	}
// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	if err := ht.bucketHeaderIO.Write(table, bucketHeader); err != nil {

// 		return err
// 	}
// 	if err := ht.objectIO.Write(table, object); err != nil {

// 		return err
// 	}
// 	return err
// }

// func (ht *TSFHashTable[T]) SaveAndSend(object T, errors chan error, mutex *sync.Mutex) {

// 	mutex.Lock()

// 	defer mutex.Unlock()

// 	table, emptyBucket, err := ht.LocateEmpty(object.Id())

// 	defer ht.fileManager.Close(table, &err)

// 	if err != nil {

// 		errors <- err

// 		return
// 	}
// 	if err := ht.fileManager.GoTo(emptyBucket, table); err != nil {

// 		errors <- err

// 		return
// 	}
// 	bucketHeader := HashTableBucketHeader{true, false, 0}

// 	if err := ht.bucketHeaderIO.Write(table, bucketHeader); err != nil {

// 		errors <- err

// 		return
// 	}
// 	if err := ht.objectIO.Write(table, object); err != nil {

// 		errors <- err

// 		return
// 	}
// }

// func (ht *TSFHashTable[T]) SaveConcurrent(objects ...T) error {

// 	errorChannel := make(chan error, len(objects))

// 	var wg sync.WaitGroup

// 	wg.Add(len(objects))

// 	var mutex sync.Mutex

// 	for _, object := range objects {

// 		go func(object T) {

// 			defer wg.Done()

// 			ht.SaveAndSend(object, errorChannel, &mutex)

// 		}(object)
// 	}
// 	go func() {

// 		wg.Wait()

// 		close(errorChannel)
// 	}()

// 	errs := make([]error, 0)

// 	for err := range errorChannel {

// 		errs = append(errs, err)
// 	}
// 	return errors.Join(errs...)
// }
