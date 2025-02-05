package schema

import "os"

type FixedSize interface {
	Size() int

	WriteSelf(file *os.File) error

	ReadSelf(file *os.File) (FixedSize, error)
}
