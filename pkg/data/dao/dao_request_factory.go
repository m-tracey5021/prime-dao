package dao

import (
	"fmt"

	"github.com/m-tracey5021/prime-dao/pkg/data"
	"github.com/m-tracey5021/prime-dao/pkg/data/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type DaoRequestFactory[T schema.Identifiable] struct {
	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	requestIds []uint64
}

func NewRequestFactory[T schema.Identifiable](fileManager fm.IFileManager, metadataManager IDaoMetadataManager[T]) DaoRequestFactory[T] {

	return DaoRequestFactory[T]{

		fileManager: fileManager,

		metadataManager: metadataManager,

		requestIds: make([]uint64, 0),
	}
}

func (factory *DaoRequestFactory[T]) NewRequestId() uint64 {

	id := data.NewId(factory.requestIds)

	factory.requestIds = append(factory.requestIds, id)

	return id
}

func (factory *DaoRequestFactory[T]) CreateSaveRequest(object T) queue.IProcessableRequest[T] {

	fmt.Printf("%v", factory.metadataManager)

	return &SaveRequest[T]{

		requestId: factory.NewRequestId(),

		object: object,

		dependencies: make(chan uint64),

		fileManager: factory.fileManager,

		metadataManager: factory.metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory *DaoRequestFactory[T]) CreateGetRequest(objectId uint64) queue.IProcessableRequest[T] {

	return &GetRequest[T]{

		requestId: factory.NewRequestId(),

		objectId: objectId,

		dependencies: make(chan uint64),

		fileManager: factory.fileManager,

		metadataManager: factory.metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory *DaoRequestFactory[T]) CreateUpdateRequest(object T) queue.IProcessableRequest[T] {

	return &UpdateRequest[T]{

		requestId: factory.NewRequestId(),

		object: object,

		dependencies: make(chan uint64),

		fileManager: factory.fileManager,

		metadataManager: factory.metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (factory *DaoRequestFactory[T]) CreateDeleteRequest(objectId uint64) queue.IProcessableRequest[T] {

	return &DeleteRequest[T]{

		requestId: factory.NewRequestId(),

		objectId: objectId,

		dependencies: make(chan uint64),

		fileManager: factory.fileManager,

		metadataManager: factory.metadataManager,

		objectIO: dataio.DataIO[T]{},
	}
}
