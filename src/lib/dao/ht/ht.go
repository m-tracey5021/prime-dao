package ht

import (
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"unsafe"
)

type ITSFHashTable[T schema.FixedSizeIdentifiable] interface {
	Save(object T) error

	QueueSaves(objects ...T)

	Get(id uint64) (*T, error)

	QueueGets(ids ...uint64)

	Update(object T) error

	Delete(id uint64) error

	GetResults() []CommandQueueResult[T]
}

type TSFHashTable[T schema.FixedSizeIdentifiable] struct {
	id uint64

	queue *CommandQueue[T]
}

func Default[T schema.FixedSizeIdentifiable](
	id uint64,

	queue *CommandQueue[T],

) TSFHashTable[T] {

	return TSFHashTable[T]{id, queue}
}

func Initialise[T schema.FixedSizeIdentifiable](
	id uint64,

	queue *CommandQueue[T],

) (ITSFHashTable[T], error) { // TODO get rid of error

	return &TSFHashTable[T]{id, queue}, nil
}

func New[T schema.FixedSizeIdentifiable](path, descriptor string, id uint64) (ITSFHashTable[T], error) {

	fileManager := fm.NewFileContainer(path, descriptor)

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(T)))

	tableSize := 10 // TODO get from config

	maxCollisions := 5 // Get from config, but should be pretty small

	hashManager := HashManager[T]{bucketSize, tableSize, maxCollisions, fileManager}

	processor := HashTableProcessor[T]{fileManager, hashManager}

	queue := NewCommandQueue(10, 5, 100, processor)

	queue.StartWorkers()

	return &TSFHashTable[T]{id, queue}, nil
}

func (ht *TSFHashTable[T]) GetResults() []CommandQueueResult[T] {

	return ht.queue.Stop()
}

func (ht *TSFHashTable[T]) Cleanup() {

	ht.queue.Stop()
}
