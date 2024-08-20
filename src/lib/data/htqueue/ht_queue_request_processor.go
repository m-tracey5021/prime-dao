package htqueue

import (
	"transformer/src/lib/data/ht"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type HashTableRequestProcessor[T schema.FixedSizeIdentifiable] struct {
	hashTable ht.TSFHashTable[T]
}

func (processor HashTableRequestProcessor[T]) Process(request queue.IRequest[T]) queue.IResult[T] {

	switch request.Type() {

	case queue.SaveRequest:

		return processor.ProcessSaveRequest(request)

	case queue.GetRequest:

		return processor.ProcessGetRequest(request)

	case queue.UpdateRequest:

		return processor.ProcessUpdateRequest(request)

	case queue.DeleteRequest:

		return processor.ProcessDeleteRequest(request)

	default:

		return nil
	}
}
