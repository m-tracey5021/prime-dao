package daoIO

import (
	"errors"
	"fmt"
	"io"
	"os"
	"transformer/src/lib/dao/schema"
)

type IFixedSizeDaoIO[T schema.FixedSize] interface {
	Write(file *os.File, object T) error

	Read(file *os.File) (T, error)

	Delete(file *os.File) error

	Zero(file *os.File) error
}

type FixedSizeDaoIO[T schema.FixedSize] struct{}

func (daoIO FixedSizeDaoIO[T]) Write(file *os.File, object T) error {

	return object.WriteSelf(file)
}

func (daoIO FixedSizeDaoIO[T]) Read(file *os.File) (T, error) {

	var object T

	read, err := object.ReadSelf(file)

	if err != nil {

		return object, err
	}
	object, ok := read.(T)

	if !ok {

		errString := fmt.Sprintf("type assertion to %T failed", object)

		return object, errors.New(errString)
	}
	return object, nil
}

func (daoIO FixedSizeDaoIO[T]) Delete(file *os.File) error {

	fileInfo, err := file.Stat()

	if err != nil {

		return err
	}
	totalLength := fileInfo.Size()

	// Find the start of the object, make sure cursor is set to this point before passing file in
	startPosition, err := file.Seek(0, io.SeekCurrent)

	if err != nil {

		return err
	}
	var object T

	endPosition := startPosition + int64(object.Size())

	initialBuffer := make([]byte, startPosition)

	remainingBuffer := make([]byte, totalLength-endPosition)

	// Go back to beginning of file
	if _, err = file.Seek(0, io.SeekStart); err != nil {

		return err
	}
	// Read everything either side of the object to be updated into a buffer
	if _, err = file.Read(initialBuffer); err != nil {

		return err
	}
	if _, err = file.Seek(endPosition, 0); err != nil {

		return err
	}
	if _, err = file.Read(remainingBuffer); err != nil {

		return err
	}
	if err = file.Truncate(0); err != nil {

		return err
	}
	if _, err = file.Seek(0, io.SeekStart); err != nil {

		return err
	}
	if _, err = file.Write(initialBuffer); err != nil {

		return err
	}
	if _, err = file.Write(remainingBuffer); err != nil {

		return err
	}
	return nil
}

func (daoIO FixedSizeDaoIO[T]) Zero(file *os.File) error {

	var object T

	size := object.Size()

	zeroBytes := make([]byte, size)

	if _, err := file.Write(zeroBytes); err != nil {

		return err
	}
	return nil
}
