package daoio

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Data int
}

func TestWrite(t *testing.T) {

	daoIO := DaoIO[TestStruct]{}

	file, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave := TestStruct{12}
	toSave2 := TestStruct{14}

	sizeWritten, err := daoIO.WriteSizePrefixed(file, toSave)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(int64(sizeWritten), 0); err != nil {

		t.Fatalf("%v", err)
	}
	_, err = daoIO.WriteSizePrefixed(file, toSave2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	read, err := readTestStruct(file)
	read2, err2 := readTestStruct(file)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err := os.Remove("test"); err != nil {

		t.Fatalf("%v", err)
	}
	if !(assert.Equal(t, toSave, *read) && assert.Equal(t, toSave2, *read2)) {

		t.Fail()
	}
}

func TestRead(t *testing.T) {

	daoIO := DaoIO[TestStruct]{}

	file, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave, sizeWritten, err := writeTestStruct(file, 12)
	toSave2, _, err2 := writeTestStruct(file, 14)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	read, err := daoIO.ReadSizePrefixed(file)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(int64(sizeWritten), 0); err != nil {

		t.Fatalf("%v", err)
	}
	read2, err := daoIO.ReadSizePrefixed(file)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err := os.Remove("test"); err != nil {

		t.Fatalf("%v", err)
	}
	if !(assert.Equal(t, toSave, read) && assert.Equal(t, toSave2, read2)) {

		t.Fail()
	}
}

func TestUpdate(t *testing.T) {

	daoIO := DaoIO[TestStruct]{}

	file, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {

		t.Fatalf("%v", err)
	}
	defer func() {

		if innerErr := os.Remove("test"); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave, sizeWritten, err := writeTestStruct(file, 12)
	_, _, err2 := writeTestStruct(file, 14)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave.Data = 43

	toSaveB := &TestStruct{51}

	if _, err := file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	sizeUpdated, err := daoIO.Update(file, *toSave)
	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(int64(sizeWritten), 0); err != nil {

		t.Fatalf("%v", err)
	}
	sizeUpdatedB, err := daoIO.Update(file, *toSaveB)

	if err != nil {

		t.Fatalf("%v", err)
	}
	fmt.Println(sizeUpdated)
	fmt.Println(sizeUpdatedB)
	if _, err := file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	read, err := readTestStruct(file)
	read2, err2 := readTestStruct(file)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err := os.Remove("test"); err != nil {

		t.Fatalf("%v", err)
	}
	if !(assert.Equal(t, toSave, read) && assert.Equal(t, toSaveB, read2)) {

		t.Fail()
	}
}

func TestDelete(t *testing.T) {

	daoIO := DaoIO[TestStruct]{}

	file, err := os.OpenFile("test", os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {

		t.Fatalf("%v", err)
	}
	_, _, errA := writeTestStruct(file, 12)
	toSaveB, sizeWrittenB, errB := writeTestStruct(file, 14)
	_, _, errC := writeTestStruct(file, 16)
	toSaveD, sizeWrittenD, _ := writeTestStruct(file, 18)

	err = errors.Join(errA, errB, errC)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	_, err = daoIO.Delete(file)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(int64(sizeWrittenB), 0); err != nil {

		t.Fatalf("%v", err)
	}
	_, err = daoIO.Delete(file)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	readA, errA := readTestStruct(file)
	readB, errB := readTestStruct(file)

	err = errors.Join(errA, errB)

	if err != nil {

		t.Fatalf("%v", err)
	}
	fileInfo, _ := file.Stat()

	if err := os.Remove("test"); err != nil {

		t.Fatalf("%v", err)
	}
	// Assert that the file is only as big as the remaining objects
	if !assert.Equal(t, sizeWrittenB+sizeWrittenD, int(fileInfo.Size())) {

		t.Fail()
	}
	if !(assert.Equal(t, toSaveB, readA) && assert.Equal(t, toSaveD, readB)) {

		t.Fail()
	}
}

func writeTestStruct(file *os.File, data int) (*TestStruct, int, error) {

	toSave := TestStruct{data}

	buffer := new(bytes.Buffer)

	encoder := gob.NewEncoder(buffer)

	if err := encoder.Encode(toSave); err != nil {

		return nil, 0, err
	}
	size := buffer.Len()

	if err := binary.Write(file, binary.LittleEndian, int64(size)); err != nil {

		return nil, 0, err
	}
	if _, err := file.Write(buffer.Bytes()); err != nil {

		return nil, 0, err
	}
	sizeWritten := size + 8

	return &toSave, sizeWritten, nil
}

func readTestStruct(file *os.File) (*TestStruct, error) {

	var size int64

	if err := binary.Read(file, binary.LittleEndian, &size); err != nil {

		return nil, err
	}
	data := make([]byte, size)

	if _, err := file.Read(data); err != nil {

		return nil, err
	}
	read := new(TestStruct)

	buffer := bytes.NewBuffer(data)

	decoder := gob.NewDecoder(buffer)

	if err := decoder.Decode(read); err != nil {

		return nil, err
	}
	return read, nil
}
