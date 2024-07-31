package dao

import (
	"fmt"
	"io"
	"os"
)

type IDaoFileContainer interface {
	ManagingFile() (*os.File, error)

	MainTable() (*os.File, error)

	CollisionTable(fileId uint64) (*os.File, error)

	GoTo(position int, file *os.File) error

	CurrentPosition(file *os.File) (int, error)

	Size(file *os.File) (int, error)

	Close(file *os.File) error
}

type DaoFileContainer struct {
	path string

	descriptor string
}

func (container DaoFileContainer) open(filePath string) (*os.File, error) {

	return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
}

func (container DaoFileContainer) ManagingFile() (*os.File, error) {

	filePath := fmt.Sprintf("%v/%v_dao", container.path, container.descriptor)

	return container.open(filePath)
}

func (container DaoFileContainer) MainTable() (*os.File, error) {

	filePath := fmt.Sprintf("%v/%v_tbl", container.path, container.descriptor)

	return container.open(filePath)
}

func (container DaoFileContainer) CollisionTable(fileId uint64) (*os.File, error) {

	filePath := fmt.Sprintf("%v/%v_c_tbl_%v", container.path, container.descriptor, fileId)

	return container.open(filePath)
}

func (container DaoFileContainer) GoTo(position int, file *os.File) error {

	if _, err := file.Seek(int64(position), io.SeekStart); err != nil {

		return err
	}
	return nil
}

func (container DaoFileContainer) CurrentPosition(file *os.File) (int, error) {

	position, err := file.Seek(0, io.SeekCurrent)

	if err != nil {

		return 0, err
	}
	return int(position), nil
}

func (container DaoFileContainer) Size(file *os.File) (int, error) {

	fileInfo, err := file.Stat()

	if err != nil {

		return 0, err
	}
	return int(fileInfo.Size()), nil
}

func (container DaoFileContainer) Close(file *os.File) error {

	return file.Close()
}
