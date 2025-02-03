package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"

	"golang.org/x/exp/maps"
)

type TSFDao[T schema.DescribedIdentifiable] struct {
	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	requestFactory *DaoRequestFactory[T]

	queue *queue.TSFQueue[T]
}

func From[T schema.DescribedIdentifiable](fileManager fm.IFileManager) (*TSFDao[T], error) {

	var described T

	fileManager = fileManager.Concatenate(described.Descriptor())

	metadataManager, err := NewMetadataManager[T](fileManager)

	if err != nil {

		return nil, err
	}
	requestFactory := NewRequestFactory(fileManager, metadataManager)

	return &TSFDao[T]{

			fileManager: fileManager,

			metadataManager: metadataManager,

			requestFactory: &requestFactory,

			queue: queue.NewQueue[T](10, 10, 100),
		},

		err
}

func (dao TSFDao[T]) NewId() uint64 {

	return dao.metadataManager.NewId()
}

func (dao TSFDao[T]) Save(object T) (int, error) {

	request := dao.requestFactory.CreateSaveRequest(object, 0)

	result := dao.queue.ProcessSync(request)

	return result.Size(), result.Error()
}

func (dao TSFDao[T]) Get(id uint64) (*T, error) {

	request := dao.requestFactory.CreateGetRequest(id, 0)

	result := dao.queue.ProcessSync(request)

	return result.Object(), result.Error()
}

func (dao TSFDao[T]) Update(object T) (int, error) {

	request := dao.requestFactory.CreateUpdateRequest(object, 0)

	result := dao.queue.ProcessSync(request)

	return result.Size(), result.Error()
}

func (dao TSFDao[T]) Delete(id uint64) (int, error) {

	request := dao.requestFactory.CreateDeleteRequest(id, 0)

	result := dao.queue.ProcessSync(request)

	return result.Size(), result.Error()
}

func (dao TSFDao[T]) AllObjectIds() []uint64 {

	return dao.metadataManager.AllObjectIds()
}

func (dao *TSFDao[T]) NewTransaction() DaoTransaction[T] {

	return DaoTransaction[T]{

		requestFactory: dao.requestFactory,

		requests: make(map[uint64]queue.IProcessableRequest[T], 0),

		dependencyResolver: queue.NewResolver[T](),

		dependencies: make([]uint64, 0),
	}
}

func (dao *TSFDao[T]) ExecuteTransaction(transaction DaoTransaction[T]) (map[uint64]queue.IResult[T], error) {

	mappedResults := make(map[uint64]queue.IResult[T])

	dao.queue = queue.NewQueue[T](10, 10, 100)

	dao.queue.Start(transaction.dependencyResolver)

	requestValues := maps.Values(transaction.requests)

	dao.queue.ProcessAsync(requestValues...)

	results := dao.queue.Stop()

	for _, result := range results {

		mappedResults[result.RequestId()] = result.Result()
	}
	if err := dao.metadataManager.SaveManagingInfo(); err != nil {

		return map[uint64]queue.IResult[T]{}, err
	}
	return mappedResults, nil
}

func (dao *TSFDao[T]) SaveMetadata() error {

	return dao.metadataManager.SaveManagingInfo()
}
