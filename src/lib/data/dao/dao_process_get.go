package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type GetRequest[T schema.Identifiable] struct {
	objectId uint64

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor GetRequest[T]) ObjectId() uint64 {

	return processor.objectId
}

func (processor GetRequest[T]) Process() queue.IResult[T] {

	index, err := processor.metadataManager.GetIndex(processor.objectId)

	if err != nil {

		return &GetResult[T]{nil, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &GetResult[T]{nil, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &GetResult[T]{nil, err}
	}
	object, err := processor.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return &GetResult[T]{nil, err}
	}
	return &GetResult[T]{object, err}
}

type GetResult[T schema.Identifiable] struct {
	object *T

	err error
}

func (result *GetResult[T]) Object() *T {

	return result.object
}

func (result *GetResult[T]) Size() int {

	return 0
}

func (result *GetResult[T]) Error() error {

	return result.err
}
