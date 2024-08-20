package dao

import (
	"transformer/src/lib/data/dataio"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/queue"
	"transformer/src/lib/data/schema"
)

type GetProcessor[T schema.Identifiable] struct {
	id uint64

	fileManager fm.IFileManager

	metadataManager *DaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

type ProcessGetResult[T schema.Identifiable] struct {
	object *T

	err error
}

func (result *ProcessGetResult[T]) Object() *T {

	return result.object
}

func (result *ProcessGetResult[T]) Size() int {

	return 0
}

func (result *ProcessGetResult[T]) Error() error {

	return result.err
}

func (processor GetProcessor[T]) Process() queue.IResult[T] {

	index, err := processor.metadataManager.GetIndex(processor.id)

	if err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, index.fileId)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	if err := processor.fileManager.GoTo(int(index.filePosition), table); err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	object, err := processor.objectIO.ReadSizePrefixed(table)

	if err != nil {

		return &ProcessGetResult[T]{nil, err}
	}
	return &ProcessGetResult[T]{object, err}
}
