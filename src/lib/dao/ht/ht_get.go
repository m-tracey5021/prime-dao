package ht

func (ht *TSFHashTable[T]) Get(id uint64) (*T, error) {

	requestId := ht.queue.NewRequestId()

	request := HashTableGetRequest[T]{requestId, id}

	result := ht.queue.Processor().Process(&request)

	return result.Object(), result.Error()
}

func (ht *TSFHashTable[T]) QueueGet(id uint64) {

	requestId := ht.queue.NewRequestId()

	request := HashTableGetRequest[T]{requestId, id}

	ht.queue.AddTask(&request)
}

func (ht *TSFHashTable[T]) QueueGets(ids ...uint64) {

	requests := make([]IRequest[T], 0)

	for _, id := range ids {

		requestId := ht.queue.NewRequestId()

		requests = append(requests, &HashTableGetRequest[T]{requestId, id})
	}
	ht.queue.AddTasks(requests)
}

// func (ht *TSFHashTable[T]) GetAndSend(id uint64, results chan *T, errors chan error, mutex *sync.Mutex) {

// 	mutex.Lock()

// 	defer mutex.Unlock()

// 	table, bucket, err := ht.Locate(id)

// 	defer ht.fileManager.Close(table, &err)

// 	if err != nil {

// 		errors <- err

// 		return
// 	}
// 	results <- &bucket.object
// }

// func (ht *TSFHashTable[T]) GetConcurrent(ids ...uint64) ([]*T, error) {

// 	objectChannel := make(chan *T, len(ids))

// 	errorChannel := make(chan error, len(ids))

// 	var wg sync.WaitGroup

// 	wg.Add(len(ids))

// 	var mutex sync.Mutex

// 	for _, id := range ids {

// 		go func(id uint64) {

// 			defer wg.Done()

// 			ht.GetAndSend(id, objectChannel, errorChannel, &mutex)

// 		}(id)
// 	}
// 	go func() {

// 		wg.Wait()

// 		close(objectChannel)

// 		close(errorChannel)
// 	}()

// 	objects := make([]*T, 0)

// 	errs := make([]error, 0)

// 	var wgResults sync.WaitGroup

// 	wgResults.Add(2)

// 	go func() {

// 		defer wgResults.Done()

// 		for object := range objectChannel {

// 			objects = append(objects, object)
// 		}
// 	}()

// 	go func() {

// 		defer wgResults.Done()

// 		for err := range errorChannel {

// 			errs = append(errs, err)
// 		}
// 	}()

// 	wgResults.Wait()

// 	return objects, errors.Join(errs...)
// }
