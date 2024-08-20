package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type DaoRequestFactory[T schema.Identifiable] struct {
}

func (factory DaoRequestFactory[T]) CreateSaveRequest(daoId uint64, fileManager fm.IFileManager, metadataManager *DaoMetadataManager[T], object T) queue.IProcessor[T] {

	return &SaveProcessor[T]{

		object: object,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory DaoRequestFactory[T]) CreateGetRequest(daoId uint64, fileManager fm.IFileManager, metadataManager *DaoMetadataManager[T], objectId uint64) queue.IProcessor[T] {

	return &GetProcessor[T]{

		id: objectId,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory DaoRequestFactory[T]) CreateUpdateRequest(daoId uint64, fileManager fm.IFileManager, metadataManager *DaoMetadataManager[T], object T) queue.IProcessor[T] {

	return &UpdateProcessor[T]{

		object: object,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory DaoRequestFactory[T]) CreateDeleteRequest(daoId uint64, fileManager fm.IFileManager, metadataManager *DaoMetadataManager[T], objectId uint64) queue.IProcessor[T] {

	return &DeleteProcessor[T]{

		id: objectId,

		fileManager: fileManager,

		metadataManager: metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}
