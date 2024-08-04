package dao

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type IFileManager interface {
	Open(fileAlias FileAlias, id uint64) (*os.File, error)

	GoTo(position int, file *os.File) error

	CurrentPosition(file *os.File) (int, error)

	Size(file *os.File) (int, error)

	Remove(file *os.File) error

	Close(file *os.File) error
}

type FileAlias int

const (
	HashTableManagingFile FileAlias = iota

	HashTableMainTable

	HashTableCollisionTable

	DaoManagingFile

	DaoTable
)

type FileManager struct {
	fileMap map[FileAlias]string
}

func NewFileContainer(path, descriptor string) FileManager {

	fileMap := map[FileAlias]string{

		HashTableManagingFile: fmt.Sprintf("%v/%v_fs_dao", path, descriptor),

		HashTableMainTable: fmt.Sprintf("%v/%v_fs_tbl", path, descriptor),

		HashTableCollisionTable: fmt.Sprintf("%v/%v_fs_c_tbl", path, descriptor),

		DaoManagingFile: fmt.Sprintf("%v/%v_dao", path, descriptor),

		DaoTable: fmt.Sprintf("%v/%v_tbl", path, descriptor),
	}
	return FileManager{fileMap}
}

func (fileManager FileManager) Open(fileAlias FileAlias, id uint64) (*os.File, error) {

	filePath, ok := fileManager.fileMap[fileAlias]

	if ok {

		filePath += fmt.Sprintf("_%v", id)

		return os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)
	}
	return nil, errors.New("file alias does not exist")
}

func (fileManager FileManager) GoTo(position int, file *os.File) error {

	if _, err := file.Seek(int64(position), io.SeekStart); err != nil {

		return err
	}
	return nil
}

func (fileManager FileManager) CurrentPosition(file *os.File) (int, error) {

	position, err := file.Seek(0, io.SeekCurrent)

	if err != nil {

		return 0, err
	}
	return int(position), nil
}

func (fileManager FileManager) Size(file *os.File) (int, error) {

	fileInfo, err := file.Stat()

	if err != nil {

		return 0, err
	}
	return int(fileInfo.Size()), nil
}

func (fileManager FileManager) Remove(file *os.File) error {

	return os.Remove(file.Name())
}

func (fileManager FileManager) Close(file *os.File) error {

	return file.Close()
}
