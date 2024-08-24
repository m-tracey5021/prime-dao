package dao

import (
	"transformer/src/lib/data"
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoRequestFactory[T schema.Identifiable] struct {
	requestIds []uint64
}

func (factory *DaoRequestFactory[T]) NewRequestId() uint64 {

	id := data.NewId(factory.requestIds)

	factory.requestIds = append(factory.requestIds, id)

	return id
}
func (factory *DaoRequestFactory[T]) CreateSaveRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], object T) queue.IProcessableRequest[T] {

	return &SaveRequest[T]{

		requestId: factory.NewRequestId(),

		object: object,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory *DaoRequestFactory[T]) CreateGetRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], objectId uint64) queue.IProcessableRequest[T] {

	return &GetRequest[T]{

		requestId: factory.NewRequestId(),

		objectId: objectId,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory *DaoRequestFactory[T]) CreateUpdateRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], object T) queue.IProcessableRequest[T] {

	return &UpdateRequest[T]{

		requestId: factory.NewRequestId(),

		object: object,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory *DaoRequestFactory[T]) CreateDeleteRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], objectId uint64) queue.IProcessableRequest[T] {

	return &DeleteRequest[T]{

		requestId: factory.NewRequestId(),

		objectId: objectId,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}
