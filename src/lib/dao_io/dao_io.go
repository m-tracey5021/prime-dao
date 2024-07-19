package daoIO

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"os"
)

type IDaoIO[T any] interface {
	GoTo(position int) error

	Write(object T) (int, error)

	WriteBool(value bool) error

	WriteAt(object T, position int) (int, error)

	WriteBoolAt(value bool, position int) error

	Read() (*T, error)

	ReadBool() (bool, error)

	ReadAt(position int) (*T, error)

	ReadBoolAt(position int) (bool, error)

	Update(object T) error

	UpdateAt(object T, position int) error

	Delete() (int, error)

	DeleteAt(position int) (int, error)

	Close() error
}

type DaoIO[T any] struct {
	file *os.File
}

func NewDaoIO[T any](filePath string) (*DaoIO[T], error) {

	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {

		return nil, err
	}
	return &DaoIO[T]{file}, nil
}

func (daoIO *DaoIO[T]) GoTo(position int) error {

	// fileInfo, err := daoIO.file.Stat()

	// if err != nil {

	// 	return err
	// }
	// size := int(fileInfo.Size())

	// if position > size {

	// 	extendedSize :=

	// 	daoIO.file.Truncate()
	// }
	_, err := daoIO.file.Seek(int64(position), io.SeekStart)

	return err
}

func (daoIO *DaoIO[T]) Write(object T) (int, error) {

	// Encode the object
	buffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(object); err != nil {

		return 0, err
	}
	// Calculate size of encoded data
	size := buffer.Len()

	// Write size as a prefix
	if err := binary.Write(daoIO.file, binary.LittleEndian, int64(size)); err != nil {

		return 0, err
	}
	// Write encoded object data
	if _, err := daoIO.file.Write(buffer.Bytes()); err != nil {

		return 0, err
	}
	// Add 8 for the size written initially
	return size + 8, nil
}

func (daoIO *DaoIO[T]) WriteBool(value bool) error {

	if err := binary.Write(daoIO.file, binary.LittleEndian, &value); err != nil {

		return err
	}
	return nil
}

func (daoIO *DaoIO[T]) WriteAt(object T, position int) (int, error) {

	if _, err := daoIO.file.Seek(int64(position), 0); err != nil {

		return 0, err
	}
	return daoIO.Write(object)
}

func (daoIO *DaoIO[T]) WriteBoolAt(value bool, position int) error {

	if _, err := daoIO.file.Seek(int64(position), 0); err != nil {

		return err
	}
	return daoIO.WriteBool(value)
}

func (daoIO *DaoIO[T]) Read() (*T, error) {

	// Read size prefix
	var size int64

	if err := binary.Read(daoIO.file, binary.LittleEndian, &size); err != nil {

		return nil, err
	}
	// Read encoded object data
	data := make([]byte, size)

	if _, err := daoIO.file.Read(data); err != nil {

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

func (daoIO *DaoIO[T]) ReadBool() (bool, error) {

	var value bool

	if err := binary.Read(daoIO.file, binary.LittleEndian, &value); err != nil {

		return false, err
	}
	return value, nil
}

func (daoIO *DaoIO[T]) ReadAt(position int) (*T, error) {

	if _, err := daoIO.file.Seek(int64(position), 0); err != nil {

		return nil, err
	}
	return daoIO.Read()
}

func (daoIO *DaoIO[T]) ReadBoolAt(position int) (bool, error) {

	if _, err := daoIO.file.Seek(int64(position), 0); err != nil {

		return false, err
	}
	return daoIO.ReadBool()
}

func (daoIO *DaoIO[T]) Update(object T) error {

	fileInfo, err := daoIO.file.Stat()

	if err != nil {

		return err
	}
	totalLength := fileInfo.Size()

	// Find the start of the object, make sure cursor is set to this point before passing file in
	startPosition, err := daoIO.file.Seek(0, io.SeekCurrent)

	if err != nil {

		return err
	}
	// Read the object to be updated
	if _, err := daoIO.Read(); err != nil {

		return err
	}
	endPosition, err := daoIO.file.Seek(0, io.SeekCurrent)

	if err != nil {

		return err
	}
	initialBuffer := make([]byte, startPosition)

	remainingBuffer := make([]byte, totalLength-endPosition)

	// Go back to beginning of file
	if _, err = daoIO.file.Seek(0, 0); err != nil {

		return err
	}
	// Read everything either side of the object to be updated into a buffer
	if _, err = daoIO.file.Read(initialBuffer); err != nil {

		return err
	}
	if _, err = daoIO.file.Seek(endPosition, 0); err != nil {

		return err
	}
	if _, err = daoIO.file.Read(remainingBuffer); err != nil {

		return err
	}
	if err = daoIO.file.Truncate(0); err != nil {

		return err
	}
	if _, err = daoIO.file.Write(initialBuffer); err != nil {

		return err
	}
	if _, err := daoIO.Write(object); err != nil {

		return err
	}
	if _, err = daoIO.file.Write(remainingBuffer); err != nil {

		return err
	}
	return nil
}

func (daoIO *DaoIO[T]) UpdateAt(object T, position int) error {

	if _, err := daoIO.file.Seek(int64(position), 0); err != nil {

		return err
	}
	return daoIO.Update(object)
}

func (daoIO *DaoIO[T]) Delete() (int, error) {

	fileInfo, err := daoIO.file.Stat()

	if err != nil {

		return 0, err
	}
	totalLength := fileInfo.Size()

	// Find the start of the object, make sure cursor is set to this point before passing file in
	startPosition, err := daoIO.file.Seek(0, io.SeekCurrent)

	if err != nil {

		return 0, err
	}
	// Read the object to be updated
	if _, err := daoIO.Read(); err != nil {

		return 0, err
	}
	endPosition, err := daoIO.file.Seek(0, io.SeekCurrent)

	sizeDeleted := endPosition - startPosition

	if err != nil {

		return 0, err
	}
	initialBuffer := make([]byte, startPosition)

	remainingBuffer := make([]byte, totalLength-endPosition)

	// Go back to beginning of file
	if _, err = daoIO.file.Seek(0, 0); err != nil {

		return 0, err
	}
	// Read everything either side of the object to be updated into a buffer
	if _, err = daoIO.file.Read(initialBuffer); err != nil {

		return 0, err
	}
	if _, err = daoIO.file.Seek(endPosition, 0); err != nil {

		return 0, err
	}
	if _, err = daoIO.file.Read(remainingBuffer); err != nil {

		return 0, err
	}
	if err = daoIO.file.Truncate(0); err != nil {

		return 0, err
	}
	if _, err = daoIO.file.Write(initialBuffer); err != nil {

		return 0, err
	}
	if _, err = daoIO.file.Write(remainingBuffer); err != nil {

		return 0, err
	}
	// Return the size that has been removed from the file
	return int(sizeDeleted), nil
}

func (daoIO *DaoIO[T]) DeleteAt(position int) (int, error) {

	if _, err := daoIO.file.Seek(int64(position), 0); err != nil {

		return 0, err
	}
	return daoIO.Delete()
}

func (daoIO *DaoIO[T]) Close() error {

	return daoIO.file.Close()
}
