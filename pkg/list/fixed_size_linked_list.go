package list

import (
	"errors"
	"io"
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type Iterator[T schema.FixedSize] struct {
	current int

	fileManager fm.IFileManager

	objectFile *os.File

	err *error

	objectIO dataio.FixedSizeDataIO[T]
}

func (list FixedSizeLinkedList[T]) Iter() (Iterator[T], error) {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	if err != nil {

		return Iterator[T]{}, err
	}
	return Iterator[T]{-1, list.fileManager, objectFile, &err, dataio.FixedSizeDataIO[T]{}}, nil
}

func (iter *Iterator[T]) Next() (*T, error) {

	read, err := iter.objectIO.Read(iter.objectFile)

	if errors.Is(err, io.EOF) {

		return nil, nil
	}
	iter.current += 1

	return &read, nil
}

func (iter *Iterator[T]) Current() int {

	return iter.current
}

func (iter Iterator[T]) Close(err *error) {

	iter.fileManager.CloseAndUnlock(iter.objectFile, iter.err)
}

type FixedSizeLinkedList[T schema.FixedSize] struct {
	id uuid.UUID

	size int

	fileManager fm.IFileManager

	objectIO dataio.FixedSizeDataIO[T]
}

func NewList[T schema.FixedSize](fileManager fm.IFileManager) FixedSizeLinkedList[T] {

	id := uuid.New()

	return FixedSizeLinkedList[T]{

		id: id,

		fileManager: fileManager,

		objectIO: dataio.FixedSizeDataIO[T]{},
	}
}

func From[T schema.FixedSize](fileManager fm.IFileManager, id uuid.UUID) FixedSizeLinkedList[T] {

	objectFile, err := fileManager.OpenAndLock(fm.ListFile, id.String())

	defer fileManager.CloseAndUnlock(objectFile, &err)

	fileSize, err := fileManager.Size(objectFile)

	var t T

	return FixedSizeLinkedList[T]{

		id: id,

		size: fileSize / t.Size(),

		fileManager: fileManager,

		objectIO: dataio.FixedSizeDataIO[T]{},
	}
}

func (list FixedSizeLinkedList[T]) Id() uuid.UUID {

	return list.id
}

func (list FixedSizeLinkedList[T]) Size() int {

	return list.size
}

func (list FixedSizeLinkedList[T]) Copy(fileManager fm.IFileManager, index int) (FixedSizeLinkedList[T], error) {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	var t T

	position := t.Size() * index

	totalLength, err := list.fileManager.Size(objectFile)

	copyLength := totalLength - position

	copyBuffer := make([]byte, copyLength)

	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return FixedSizeLinkedList[T]{}, err
	}
	if _, err = objectFile.Read(copyBuffer); err != nil {

		return FixedSizeLinkedList[T]{}, err
	}
	newId := uuid.New()

	newObjectFile, err := list.fileManager.OpenAndLock(fm.ListFile, newId.String())

	defer list.fileManager.CloseAndUnlock(newObjectFile, &err)

	if _, err = newObjectFile.Write(copyBuffer); err != nil {

		return FixedSizeLinkedList[T]{}, err
	}
	return FixedSizeLinkedList[T]{

		id: newId,

		size: copyLength / t.Size(),

		fileManager: fileManager,

		objectIO: dataio.FixedSizeDataIO[T]{},
	}, err
}

func (list FixedSizeLinkedList[T]) Replace(object T, index int) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	var t T

	position := index * t.Size()

	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return err
	}
	err = list.objectIO.Write(objectFile, object)

	return err
}

// TODO sort out errors being covered by return statement
func (list FixedSizeLinkedList[T]) Index(index int) (T, error) {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	var t T

	position := index * t.Size()

	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return t, err
	}
	if object, err := list.objectIO.Read(objectFile); err != nil {

		return t, err

	} else {

		return object, err
	}
}

func (list FixedSizeLinkedList[T]) ToSlice() ([]T, error) {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	slice := []T{}

	for range list.size {

		object, err := list.objectIO.Read(objectFile)

		if err != nil {

			return []T{}, err
		}
		slice = append(slice, object)
	}
	return slice, err
}

func (list *FixedSizeLinkedList[T]) Insert(object T, index int) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	var t T

	position := index * t.Size()

	totalLength, err := list.fileManager.Size(objectFile)

	if err != nil {

		return err
	}
	remainingBuffer := make([]byte, totalLength-position)

	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return err
	}
	if _, err = objectFile.Read(remainingBuffer); err != nil {

		return err
	}
	if err = objectFile.Truncate(int64(position)); err != nil {

		return err
	}
	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return err
	}
	if err := list.objectIO.Write(objectFile, object); err != nil {

		return err
	}
	_, err = objectFile.Write(remainingBuffer)

	list.size += 1

	return err
}

func (list *FixedSizeLinkedList[T]) Append(object T) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	size, err := list.fileManager.Size(objectFile)

	if err != nil {

		return err
	}
	if err := list.fileManager.GoTo(size, objectFile); err != nil {

		return err
	}
	err = list.objectIO.Write(objectFile, object)

	list.size += 1

	return err
}

func (list *FixedSizeLinkedList[T]) AppendAll(other FixedSizeLinkedList[T]) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	toAppend, err := list.fileManager.OpenAndLock(fm.ListFile, other.id.String())

	defer list.fileManager.CloseAndUnlock(toAppend, &err)

	initialSize, err := list.fileManager.Size(objectFile)

	if err != nil {

		return err
	}
	sizeToAppend, err := list.fileManager.Size(toAppend)

	if err != nil {

		return err
	}
	buffer := make([]byte, sizeToAppend)

	if _, err := toAppend.Read(buffer); err != nil {

		return err
	}
	if err := list.fileManager.GoTo(initialSize, objectFile); err != nil {

		return err
	}
	if _, err = objectFile.Write(buffer); err != nil {

		return err
	}
	list.size += other.Size()

	return err
}

func (list *FixedSizeLinkedList[T]) Remove(index int) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	var t T

	position := index * t.Size()

	start := (index + 1) * t.Size()

	totalLength, err := list.fileManager.Size(objectFile)

	if err != nil {

		return err
	}
	remainingBuffer := make([]byte, totalLength-start)

	if err := list.fileManager.GoTo(start, objectFile); err != nil {

		return err
	}
	if _, err = objectFile.Read(remainingBuffer); err != nil {

		return err
	}
	if err = objectFile.Truncate(int64(position)); err != nil {

		return err
	}
	if err := list.fileManager.GoTo(position, objectFile); err != nil {

		return err
	}
	_, err = objectFile.Write(remainingBuffer)

	list.size -= 1

	return err
}

func (list *FixedSizeLinkedList[T]) Truncate(index int) error {

	objectFile, err := list.fileManager.OpenAndLock(fm.ListFile, list.id.String())

	defer list.fileManager.CloseAndUnlock(objectFile, &err)

	var t T

	position := index * t.Size()

	totalLength, err := list.fileManager.Size(objectFile)

	truncatedLength := totalLength - position

	numberOfObjectsTruncated := truncatedLength / t.Size()

	err = objectFile.Truncate(int64(position))

	list.size -= numberOfObjectsTruncated

	return err
}
