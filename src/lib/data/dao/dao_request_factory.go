package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoRequestFactory[T schema.Identifiable] struct {
}

func (factory DaoRequestFactory[T]) CreateSaveRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], object T) queue.IProcessableRequest[T] {

	return &SaveRequest[T]{

		object: object,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory DaoRequestFactory[T]) CreateGetRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], objectId uint64) queue.IProcessableRequest[T] {

	return &GetRequest[T]{

		objectId: objectId,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory DaoRequestFactory[T]) CreateUpdateRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], object T) queue.IProcessableRequest[T] {

	return &UpdateRequest[T]{

		object: object,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory DaoRequestFactory[T]) CreateDeleteRequest(daoId uint64, fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T], objectId uint64) queue.IProcessableRequest[T] {

	return &DeleteRequest[T]{

		objectId: objectId,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}
