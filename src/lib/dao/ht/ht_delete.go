package ht

func (ht *TSFHashTable[T]) Delete(id uint64) error {

	requestId := ht.queue.NewRequestId()

	request := HashTableDeleteRequest[T]{requestId, id}

	result := ht.queue.Processor().Process(&request)

	return result.Error()
}

func (ht *TSFHashTable[T]) QueueDelete(id uint64) {

	requestId := ht.queue.NewRequestId()

	request := HashTableDeleteRequest[T]{requestId, id}

	ht.queue.AddTask(&request)
}

func (ht *TSFHashTable[T]) QueueDeletes(ids ...uint64) {

	requests := make([]IRequest[T], 0)

	for _, id := range ids {

		requestId := ht.queue.NewRequestId()

		requests = append(requests, &HashTableDeleteRequest[T]{requestId, id})
	}
	ht.queue.AddTasks(requests)
}
