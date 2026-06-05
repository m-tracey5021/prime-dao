package list

import (
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type LinkedList[T schema.Identifiable] struct {
	id uuid.UUID

	size int

	fileManager fm.IFileManager

	objectIO dataio.DataIO[T]
}

func NewIdentifiableList[T schema.Identifiable](fileManager fm.IFileManager) LinkedList[T] {

	id := uuid.New()

	return LinkedList[T]{

		id: id,

		fileManager: fileManager,

		objectIO: dataio.DataIO[T]{},
	}
}

func (list LinkedList[T]) Id() uuid.UUID {

	return list.id
}

func (list LinkedList[T]) Size() int {

	return list.size
}

func (list LinkedList[T]) GetPositionOfIndex(file *os.File, index int) (int, int, error) {

	for range index - 1 {

		if _, err := list.objectIO.ReadSizePrefixed(file); err != nil {

			return 0, 0, err
		}
	}
	start, err := list.fileManager.CurrentPosition(file)

	if err != nil {

		return 0, 0, err
	}
	if _, err := list.objectIO.ReadSizePrefixed(file); err != nil {

		return 0, 0, err
	}
	end, err := list.fileManager.CurrentPosition(file)

	if err != nil {

		return 0, 0, err
	}
	return start, end, nil
}

func (list LinkedList[T]) Copy(index int) (LinkedList[T], error) {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	position, _, err := list.GetPositionOfIndex(objectFile, index)

	if err != nil {

		return LinkedList[T]{}, err
	}
	totalLength, err := list.fileManager.Size(objectFile)

	copyLength := totalLength - position

	copyBuffer := make([]byte, copyLength)

	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return LinkedList[T]{}, err
	}
	if _, err = objectFile.Read(copyBuffer); err != nil {

		return LinkedList[T]{}, err
	}
	newId := uuid.New()

	newObjectFile, err := list.fileManager.OpenAndLock(fm.ListFile, newId.String())

	defer list.fileManager.CloseAndUnlock(newObjectFile, &err)

	if _, err = newObjectFile.Write(copyBuffer); err != nil {

		return LinkedList[T]{}, err
	}
	return LinkedList[T]{

		id: newId,

		size: list.size - index,

		fileManager: list.fileManager, // have changed this from being passed in, not sure if it matters or not

		objectIO: dataio.DataIO[T]{},
	}, err
}

func (list LinkedList[T]) Replace(object T, index int) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	start, end, err := list.GetPositionOfIndex(objectFile, index)

	totalLength, err := list.fileManager.Size(objectFile)

	copyLength := totalLength - end

	copyBuffer := make([]byte, copyLength)

	if err := list.fileManager.GoTo(end, objectFile); err != nil {

		return err
	}
	if _, err = objectFile.Read(copyBuffer); err != nil {

		return err
	}
	if err := list.fileManager.GoTo(start, objectFile); err != nil {

		return err
	}
	if _, err := list.objectIO.WriteSizePrefixed(objectFile, object); err != nil {

		return err
	}
	if _, err = objectFile.Write(copyBuffer); err != nil {

		return err
	}
	return nil
}

// func (list LinkedList[T]) Index(index int) (T, error) {

// 	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

// 	defer list.fileManager.CloseAndUnlock(objectFile, &err)

// 	var t T

// 	position := index * t.Size()

// 	if err := list.fileManager.GoTo(position, objectFile); err != nil {

// 		return t, err
// 	}
// 	if object, err := list.objectIO.Read(objectFile); err != nil {

// 		return t, err

// 	} else {

// 		return object, err
// 	}
// }
