package dao

import (
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type TSFDao[T schema.Identifiable] struct {
	id uint64

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	requestFactory DaoRequestFactory[T]

	requests []queue.IProcessableRequest[T]

	queue *queue.TSFQueue[T]
}

func New[T schema.Identifiable](path, descriptor string, id uint64) (*TSFDao[T], error) {

	fileManager := fm.NewFileManager(path, descriptor)

	metadataManager, err := NewMetadataManager[T](id, fileManager)

	if err != nil {

		return nil, err
	}
	return &TSFDao[T]{

			id: id,

			fileManager: fileManager,

			metadataManager: metadataManager,

			requestFactory: DaoRequestFactory[T]{},

			requests: make([]queue.IProcessableRequest[T], 0),

			queue: &queue.TSFQueue[T]{},
		},

		err
}

func (dao *TSFDao[T]) QueueSaveRequest(request SaveRequest[T]) {

	dao.requests = append(dao.requests, request)
}

func (dao *TSFDao[T]) QueueGetRequest(request GetRequest[T]) {

	dao.requests = append(dao.requests, request)
}

func (dao *TSFDao[T]) QueueUpdateRequest(request UpdateRequest[T]) {

	dao.requests = append(dao.requests, request)
}

func (dao *TSFDao[T]) QueueDeleteRequest(request DeleteRequest[T]) {

	dao.requests = append(dao.requests, request)
}

func (dao *TSFDao[T]) GroupRequests() map[uint64][]queue.IProcessableRequest[T] {

	// make sure each group doesnt include the same Id twice
	requestMapping := make(map[uint64][]queue.IProcessableRequest[T], 0)

	for _, request := range dao.requests {

		groupedRequests, ok := requestMapping[request.ObjectId()]

		if ok {

			requestMapping[request.ObjectId()] = append(groupedRequests, request)

		} else {

			requestMapping[request.ObjectId()] = []queue.IProcessableRequest[T]{request}
		}
	}
	return requestMapping
}

func (dao *TSFDao[T]) Execute() {

	reqGroups := dao.GroupRequests()

	for _, group := range reqGroups {

		dao.queue.ProcessAsync(group...)

		// wait here
	}
}
