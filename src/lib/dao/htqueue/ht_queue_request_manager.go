package htqueue

import (
	"transformer/src/lib/dao/ht"
	"transformer/src/lib/dao/queue"
	"transformer/src/lib/dao/schema"
)

func NewRequestManager[T schema.FixedSizeIdentifiable](path, descriptor string, id uint64) *queue.QueueRequestManager[T] {

	hashTable := ht.New[T](path, descriptor, id)

	processor := HashTableRequestProcessor[T]{hashTable}

	requestQueue := queue.NewQueue(10, 5, 100, processor)

	requestFactory := HashTableRequestFactory[T]{}

	requestManager := queue.NewQueueRequestManager[T](requestFactory, requestQueue)

	return &requestManager
}
