package daoio

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
	"unsafe"
)

type IDaoIO[T any] interface {
	WriteSizePrefixed(file *os.File, object T) (int, error)

	ReadSizePrefixed(file *os.File) (*T, error)

	Update(file *os.File, object T) (int, error)

	Delete(file *os.File) (int, error)

	Zero(file *os.File) error

	Size(file *os.File) (int, error) // TODO delete this, duplicate from fileContainer
}

type DaoIO[T any] struct{}

func (daoIO DaoIO[T]) WriteSizePrefixed(file *os.File, object T) (int, error) {

	// Encode the object
	buffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(object); err != nil {

		return 0, err
	}
	// Calculate size of encoded data
	size := buffer.Len()

	// Write size as a prefix
	if err := binary.Write(file, binary.LittleEndian, uint64(size)); err != nil {

		return 0, err
	}
	// Write encoded object data
	if _, err := file.Write(buffer.Bytes()); err != nil {

		return 0, err
	}
	// Add 8 for the size written initially
	return size + 8, nil
}

func (daoIO DaoIO[T]) ReadSizePrefixed(file *os.File) (*T, error) {

	// Read size prefix
	var size uint64

	if err := binary.Read(file, binary.LittleEndian, &size); err != nil {

		return nil, err
	}
	// Read encoded object data
	data := make([]byte, size)

	if _, err := file.Read(data); err != nil {

		return nil, err
	}
	// Decode the object
	object := new(T)

	buffer := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(object); err != nil {

		return nil, err
	}
	return object, nil
}

func (daoIO DaoIO[T]) Update(file *os.File, object T) (int, error) {

	fileInfo, err := file.Stat()

	if err != nil {

		return 0, err
	}
	totalLength := fileInfo.Size()

	// Find the start of the object, make sure cursor is set to this point before passing file in
	startPosition, err := file.Seek(0, io.SeekCurrent)

	if err != nil {

		return 0, err
	}
	// Read the object to be updated
	if _, err := daoIO.ReadSizePrefixed(file); err != nil {

		return 0, err
	}
	endPosition, err := file.Seek(0, io.SeekCurrent)

	if err != nil {

		return 0, err
	}
	initialBuffer := make([]byte, startPosition)

	remainingBuffer := make([]byte, totalLength-endPosition)

	// Go back to beginning of file
	if _, err = file.Seek(0, 0); err != nil {

		return 0, err
	}
	// Read everything either side of the object to be updated into a buffer
	if _, err = file.Read(initialBuffer); err != nil {

		return 0, err
	}
	if _, err = file.Seek(endPosition, 0); err != nil {

		return 0, err
	}
	if _, err = file.Read(remainingBuffer); err != nil {

		return 0, err
	}
	if err = file.Truncate(0); err != nil {

		return 0, err
	}
	if _, err = file.Seek(0, 0); err != nil {

		return 0, err
	}
	if _, err = file.Write(initialBuffer); err != nil {

		return 0, err
	}
	if _, err := daoIO.WriteSizePrefixed(file, object); err != nil {

		return 0, err
	}
	if _, err = file.Write(remainingBuffer); err != nil {

		return 0, err
	}
	updatedLength := fileInfo.Size()

	return int(updatedLength - totalLength), nil
}

func (daoIO DaoIO[T]) Delete(file *os.File) (int, error) {

	fileInfo, err := file.Stat()

	if err != nil {

		return 0, err
	}
	totalLength := fileInfo.Size()

	// Find the start of the object, make sure cursor is set to this point before passing file in
	startPosition, err := file.Seek(0, io.SeekCurrent)

	if err != nil {

		return 0, err
	}
	// Read the object to be updated
	if _, err := daoIO.ReadSizePrefixed(file); err != nil {

		return 0, err
	}
	endPosition, err := file.Seek(0, io.SeekCurrent)

	sizeDeleted := endPosition - startPosition // NOT RIGHT

	if err != nil {

		return 0, err
	}
	initialBuffer := make([]byte, startPosition)

	remainingBuffer := make([]byte, totalLength-endPosition)

	// Go back to beginning of file
	if _, err = file.Seek(0, 0); err != nil {

		return 0, err
	}
	// Read everything either side of the object to be updated into a buffer
	if _, err = file.Read(initialBuffer); err != nil {

		return 0, err
	}
	if _, err = file.Seek(endPosition, 0); err != nil {

		return 0, err
	}
	if _, err = file.Read(remainingBuffer); err != nil {

		return 0, err
	}
	if err = file.Truncate(0); err != nil {

		return 0, err
	}
	if _, err = file.Seek(0, 0); err != nil {

		return 0, err
	}
	if _, err = file.Write(initialBuffer); err != nil {

		return 0, err
	}
	if _, err = file.Write(remainingBuffer); err != nil {

		return 0, err
	}
	// Return the size that has been removed from the file
	return int(sizeDeleted), nil
}

func (daoIO DaoIO[T]) Zero(file *os.File) error {

	size := int(unsafe.Sizeof(*new(T)))

	zeroBuffer := make([]byte, size)

	_, err := file.Write(zeroBuffer)

	if err != nil {

		return err
	}
	return nil
}

func (daoIO DaoIO[T]) Size(file *os.File) (int, error) {

	fileInfo, err := file.Stat()

	if err != nil {

		return 0, err
	}
	return int(fileInfo.Size()), nil
}
