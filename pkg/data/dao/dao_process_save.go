package dao

import (
	"github.com/m-tracey5021/prime-dao/pkg/data/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/data/fm"
	"github.com/m-tracey5021/prime-dao/pkg/data/queue"
	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type SaveRequest[T schema.Identifiable] struct {
	requestId uint64

	object T

	dependencies chan uint64

	numberOfDependencies int

	fileManager fm.IFileManager

	metadataManager IDaoMetadataManager[T]

	objectIO dataio.DataIO[T]
}

func (processor *SaveRequest[T]) RequestId() uint64 {

	return processor.requestId
}

func (processor *SaveRequest[T]) ObjectId() uint64 {

	return processor.object.Id()
}

func (processor *SaveRequest[T]) Dependencies() chan uint64 {

	return processor.dependencies
}

func (processor *SaveRequest[T]) GetNumDeps() int {

	return processor.numberOfDependencies
}

func (processor *SaveRequest[T]) SetNumDeps(deps int) {

	processor.numberOfDependencies = deps
}

func (processor *SaveRequest[T]) Process() queue.IResult[T] {

	availableTable := processor.metadataManager.UpdateMetadataPreSave()

	table, err := processor.fileManager.OpenAndLock(fm.DaoTable, availableTable)

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &SaveResult[T]{0, err}
	}
	position, err := processor.fileManager.Size(table)

	if err := processor.fileManager.GoTo(position, table); err != nil {

		return &SaveResult[T]{0, err}
	}
	size, err := processor.objectIO.WriteSizePrefixed(table, processor.object)

	if err != nil {

		return &SaveResult[T]{0, err}
	}
	processor.metadataManager.UpdateMetadataPostSave(availableTable, processor.object, uint64(position))

	return &SaveResult[T]{size, err}
}

type SaveResult[T schema.Identifiable] struct {
	sizeWritten int

	err error
}

func (result *SaveResult[T]) Object() *T {

	return nil
}

func (result *SaveResult[T]) Size() int {

	return result.sizeWritten
}

func (result *SaveResult[T]) Error() error {

	return result.err
}
