package ht

func (ht *TSFHashTable[T]) Update(object T) error {

	requestId := ht.queue.NewRequestId()

	request := HashTableUpdateRequest[T]{requestId, object}

	return ht.queue.Processor().Process(&request).Error()
}

func (ht *TSFHashTable[T]) QueueUpdate(object T) {

	requestId := ht.queue.NewRequestId()

	request := HashTableUpdateRequest[T]{requestId, object}

	ht.queue.AddTask(&request)
}

func (ht *TSFHashTable[T]) QueueUpdates(objects ...T) {

	requests := make([]IRequest[T], 0)

	for _, object := range objects {

		requestId := ht.queue.NewRequestId()

		requests = append(requests, &HashTableUpdateRequest[T]{requestId, object})
	}
	ht.queue.AddTasks(requests)
}
