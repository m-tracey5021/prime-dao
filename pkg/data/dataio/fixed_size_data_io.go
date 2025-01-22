package dataio

import (
	"errors"
	"fmt"
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/data/schema"
)

type IFixedSizeDataIO[T schema.FixedSize] interface {
	Write(file *os.File, object T) error

	Read(file *os.File) (T, error)

	Zero(file *os.File) error
}

type FixedSizeDataIO[T schema.FixedSize] struct{}

func (dataIO FixedSizeDataIO[T]) Write(file *os.File, object T) error {

	return object.WriteSelf(file)
}

func Write[T schema.FixedSize](file *os.File, object T) error {

	return object.WriteSelf(file)
}

func (dataIO FixedSizeDataIO[T]) Read(file *os.File) (T, error) {

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

func (dataIO FixedSizeDataIO[T]) Zero(file *os.File) error {

	var object T

	size := object.Size()

	zeroBytes := make([]byte, size)

	if _, err := file.Write(zeroBytes); err != nil {

		return err
	}
	return nil
}
