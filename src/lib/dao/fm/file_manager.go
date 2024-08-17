package fm

import (
	"errors"
	"fmt"
	"io"
	"os"
	"syscall"
)

type IFileManager interface {
	OpenAndLock(fileAlias FileAlias, idChain ...uint64) (*os.File, error)

	GoTo(position int, file *os.File) error

	CurrentPosition(file *os.File) (int, error)

	Size(file *os.File) (int, error)

	Remove(file *os.File) error

	CloseAndUnlock(file *os.File, err *error) error
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

		HashTableManagingFile: fmt.Sprintf("%v/%v_ht", path, descriptor),

		HashTableMainTable: fmt.Sprintf("%v/%v_ht_tbl", path, descriptor),

		HashTableCollisionTable: fmt.Sprintf("%v/%v_ht_c_tbl", path, descriptor),

		DaoManagingFile: fmt.Sprintf("%v/%v_dao", path, descriptor),

		DaoTable: fmt.Sprintf("%v/%v_dao_tbl", path, descriptor),
	}
	return FileManager{fileMap}
}

func (fileManager FileManager) OpenAndLock(fileAlias FileAlias, idChain ...uint64) (*os.File, error) {

	filePath, ok := fileManager.fileMap[fileAlias]

	if ok {

		for _, id := range idChain {

			filePath += fmt.Sprintf("_%v", id)
		}
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0666)

		if err != nil {

			return nil, err
		}
		err = fileManager.Lock(file) // This line will block until the file is available

		if err != nil {

			return nil, err
		}
		return file, err
	}
	return nil, errors.New("file alias does not exist")
}

func (fileManager FileManager) Lock(file *os.File) error {

	return syscall.Flock(int(file.Fd()), syscall.LOCK_EX)
}

func (fileManager FileManager) Unlock(file *os.File) error {

	return syscall.Flock(int(file.Fd()), syscall.LOCK_UN)
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

func (fileManager FileManager) CloseAndUnlock(file *os.File, err *error) error {

	if unlockErr := fileManager.Unlock(file); unlockErr != nil {

		if *err != nil {

			*err = errors.Join(unlockErr, *err)

		} else {

			*err = unlockErr
		}
	}
	if closeErr := file.Close(); closeErr != nil {

		if *err != nil {

			*err = errors.Join(closeErr, *err)

		} else {

			*err = closeErr
		}
	}
	return *err
}

func (fileManager FileManager) Remove(file *os.File) error {

	return os.Remove(file.Name())
}
