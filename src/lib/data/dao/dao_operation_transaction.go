package dao

import (
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoTransaction[T schema.Identifiable] struct {
	requests []queue.IProcessableRequest[T]
}

func (transaction *DaoTransaction[T]) AddRequest(request queue.IProcessableRequest[T]) {

	transaction.requests = append(transaction.requests, request)
}

func (transaction *DaoTransaction[T]) GroupRequests() map[uint64][]queue.IProcessableRequest[T] {

	requestMapping := make(map[uint64][]queue.IProcessableRequest[T], 0)

	for _, request := range transaction.requests {

		groupedRequests, ok := requestMapping[request.ObjectId()]

		if ok {

			requestMapping[request.ObjectId()] = append(groupedRequests, request)

		} else {

			requestMapping[request.ObjectId()] = []queue.IProcessableRequest[T]{request}
		}
	}
	return requestMapping
}
