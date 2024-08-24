package dao

import (
	"slices"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type TSFDao[T schema.Identifiable] struct {
	id uint64

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	requestFactory *DaoRequestFactory[T]

	requests []queue.IProcessableRequest[T]

	queueContext *queue.QueueContext[T]

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

			requestFactory: &DaoRequestFactory[T]{},

			requests: make([]queue.IProcessableRequest[T], 0),

			queueContext: queue.NewContext[T](),

			queue: queue.NewQueue[T](10, 10, 100),
		},

		err
}

func (dao *TSFDao[T]) QueueSaveRequest(object T, dependsOn ...uint64) uint64 {

	request := dao.requestFactory.CreateSaveRequest(dao.id, dao.fileManager, dao.metadataManager, object)

	dao.requests = append(dao.requests, request)

	dao.queueContext.AddDependency(request, dependsOn)

	return request.RequestId()
}

func (dao *TSFDao[T]) QueueGetRequest(id uint64, dependsOn ...uint64) uint64 {

	request := dao.requestFactory.CreateGetRequest(dao.id, dao.fileManager, dao.metadataManager, id, len(dependsOn))

	dao.requests = append(dao.requests, request)

	dao.queueContext.AddDependency(request, dependsOn)

	return request.RequestId()
}

func (dao *TSFDao[T]) QueueUpdateRequest(object T, dependsOn ...uint64) uint64 {

	request := dao.requestFactory.CreateUpdateRequest(dao.id, dao.fileManager, dao.metadataManager, object)

	dao.requests = append(dao.requests, request)

	dao.queueContext.AddDependency(request, dependsOn)

	return request.RequestId()
}

func (dao *TSFDao[T]) QueueDeleteRequest(id uint64, dependsOn ...uint64) uint64 {

	request := dao.requestFactory.CreateDeleteRequest(dao.id, dao.fileManager, dao.metadataManager, id)

	dao.requests = append(dao.requests, request)

	dao.queueContext.AddDependency(request, dependsOn)

	return request.RequestId()
}

func (dao *TSFDao[T]) GroupRequests() [][]queue.IProcessableRequest[T] {

	groups := make([][]queue.IProcessableRequest[T], 0)

	for _, request := range dao.requests {

		added := false

		for _, group := range groups {

			found := slices.ContainsFunc(group, func(element queue.IProcessableRequest[T]) bool {

				return element.ObjectId() == request.ObjectId()
			})
			if !found {

				added = true

				group[request.ObjectId()] = request
			}
		}
		if !added {

			groups = append(groups, []queue.IProcessableRequest[T]{request})
		}
	}
	return groups
}

func (dao *TSFDao[T]) Execute() map[uint64]queue.IResult[T] {

	mappedResults := make(map[uint64]queue.IResult[T])

	reqGroups := dao.GroupRequests()

	for _, group := range reqGroups {

		dao.queue = queue.NewQueue[T](10, 10, 100)

		dao.queue.Start(dao.queueContext)

		dao.queue.ProcessAsync(group...)

		results := dao.queue.Stop()

		for _, result := range results {

			mappedResults[result.RequestId()] = result.Result()
		}
	}
	dao.requests = make([]queue.IProcessableRequest[T], 0)

	return mappedResults
}

func (dao *TSFDao[T]) ExecuteReq() map[uint64]queue.IResult[T] {

	mappedResults := make(map[uint64]queue.IResult[T])

	dao.queue = queue.NewQueue[T](10, 10, 100)

	dao.queue.Start(dao.queueContext)

	dao.queue.ProcessAsync(dao.requests...)

	results := dao.queue.Stop()

	for _, result := range results {

		mappedResults[result.RequestId()] = result.Result()
	}
	dao.requests = make([]queue.IProcessableRequest[T], 0)

	return mappedResults
}
