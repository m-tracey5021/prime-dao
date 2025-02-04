package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type Dao[T schema.DescribedIdentifiable] struct {
	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	requestFactory *DaoRequestFactory[T]

	queue *queue.Queue[T]
}

func From[T schema.DescribedIdentifiable](fileManager fm.IFileManager) (*Dao[T], error) {

	var described T

	fileManager = fileManager.Concatenate(described.Descriptor())

	metadataManager, err := NewMetadataManager[T](fileManager)

	if err != nil {

		return nil, err
	}
	requestFactory := NewRequestFactory(fileManager, metadataManager)

	return &Dao[T]{

			fileManager: fileManager,

			metadataManager: metadataManager,

			requestFactory: &requestFactory,

			queue: queue.NewQueue[T](10, 10, 100),
		},

		err
}

func (dao Dao[T]) NewId() uint64 {

	return dao.metadataManager.NewId()
}

func (dao Dao[T]) Save(object T) (int, error) {

	request := dao.requestFactory.CreateSaveRequest(object)

	result := dao.queue.ProcessSync(request)

	return result.Size(), result.Error()
}

func (dao Dao[T]) Get(id uint64) (*T, error) {

	request := dao.requestFactory.CreateGetRequest(id)

	result := dao.queue.ProcessSync(request)

	return result.Object(), result.Error()
}

func (dao Dao[T]) Update(object T) (int, error) {

	request := dao.requestFactory.CreateUpdateRequest(object)

	result := dao.queue.ProcessSync(request)

	return result.Size(), result.Error()
}

func (dao Dao[T]) Delete(id uint64) (int, error) {

	request := dao.requestFactory.CreateDeleteRequest(id)

	result := dao.queue.ProcessSync(request)

	return result.Size(), result.Error()
}

func (dao Dao[T]) AllObjectIds() []uint64 {

	return dao.metadataManager.AllObjectIds()
}

func (dao *Dao[T]) NewTransaction() DaoTransaction[T] {

	return DaoTransaction[T]{

		requestFactory: dao.requestFactory,

		requests: []queue.IProcessableRequest[T]{},

		dependencyResolver: queue.NewResolver[T](),

		dependencies: make([]uint64, 0),
	}
}

func (dao *Dao[T]) ExecuteTransaction(transaction DaoTransaction[T]) (map[uint64]queue.IResult[T], error) {

	mappedResults := make(map[uint64]queue.IResult[T])

	dao.queue = queue.NewQueue[T](10, 10, 100)

	dao.queue.Start(transaction.dependencyResolver)

	dao.queue.ProcessAsync(transaction.requests...)

	results := dao.queue.Stop()

	for _, result := range results {

		mappedResults[result.RequestId()] = result.Result()
	}
	if err := dao.metadataManager.SaveManagingInfo(); err != nil {

		return map[uint64]queue.IResult[T]{}, err
	}
	return mappedResults, nil
}
