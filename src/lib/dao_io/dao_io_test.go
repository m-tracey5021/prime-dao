package daoIO

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type TestStruct struct {
	Data int
}

func TestWrite(t *testing.T) {

	daoIO, err := NewDaoIO[TestStruct]("test.gob")

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave := TestStruct{12}
	toSave2 := TestStruct{14}

	sizeWritten, err := daoIO.Write(toSave)

	if err != nil {

		t.Fatalf("%v", err)
	}
	_, err = daoIO.WriteAt(toSave2, sizeWritten)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := daoIO.file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	read, err := readTestStruct(daoIO.file)
	read2, err2 := readTestStruct(daoIO.file)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err := os.Remove("test.gob"); err != nil {

		t.Fatalf("%v", err)
	}
	if !(assert.Equal(t, toSave, *read) && assert.Equal(t, toSave2, *read2)) {

		t.Fail()
	}
}

func TestRead(t *testing.T) {

	daoIO, err := NewDaoIO[TestStruct]("test.gob")

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave, sizeWritten, err := writeTestStruct(daoIO.file, 12)
	toSave2, _, err2 := writeTestStruct(daoIO.file, 14)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := daoIO.file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	read, err := daoIO.Read()

	if err != nil {

		t.Fatalf("%v", err)
	}
	read2, err := daoIO.ReadAt(sizeWritten)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err := os.Remove("test.gob"); err != nil {

		t.Fatalf("%v", err)
	}
	if !(assert.Equal(t, toSave, read) && assert.Equal(t, toSave2, read2)) {

		t.Fail()
	}
}

func TestUpdate(t *testing.T) {

	daoIO, err := NewDaoIO[TestStruct]("test.gob")

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave, sizeWritten, err := writeTestStruct(daoIO.file, 12)
	toSave2, _, err2 := writeTestStruct(daoIO.file, 14)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	toSave.Data = 43

	toSave2 = &TestStruct{51}

	if _, err := daoIO.file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	daoIO.Update(*toSave)

	daoIO.UpdateAt(*toSave2, sizeWritten)

	if _, err := daoIO.file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	read, err := readTestStruct(daoIO.file)
	read2, err2 := readTestStruct(daoIO.file)

	err = errors.Join(err, err2)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if err := os.Remove("test.gob"); err != nil {

		t.Fatalf("%v", err)
	}
	if !(assert.Equal(t, toSave, read) && assert.Equal(t, toSave2, read2)) {

		t.Fail()
	}
}

func TestDelete(t *testing.T) {

	daoIO, err := NewDaoIO[TestStruct]("test.gob")

	if err != nil {

		t.Fatalf("%v", err)
	}
	_, _, errA := writeTestStruct(daoIO.file, 12)
	toSaveB, sizeWrittenB, errB := writeTestStruct(daoIO.file, 14)
	_, _, errC := writeTestStruct(daoIO.file, 16)
	toSaveD, sizeWrittenD, _ := writeTestStruct(daoIO.file, 18)

	err = errors.Join(errA, errB, errC)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := daoIO.file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	_, err = daoIO.Delete()

	if err != nil {

		t.Fatalf("%v", err)
	}
	_, err = daoIO.DeleteAt(sizeWrittenB)

	if err != nil {

		t.Fatalf("%v", err)
	}
	if _, err := daoIO.file.Seek(0, 0); err != nil {

		t.Fatalf("%v", err)
	}
	readA, errA := readTestStruct(daoIO.file)
	readB, errB := readTestStruct(daoIO.file)

	err = errors.Join(errA, errB)

	if err != nil {

		t.Fatalf("%v", err)
	}
	fileInfo, _ := daoIO.file.Stat()

	if err := os.Remove("test.gob"); err != nil {

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
