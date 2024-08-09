package dao_test

import (
	"testing"
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/ht"

	"github.com/stretchr/testify/assert"
)

func TestHashableDao(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	fileManager := fm.NewFileContainer(path, descriptor)

	daoId := uint64(0)

	hashableDao, err := ht.New[ht.MockHashable](fileManager, daoId)

	if err != nil {

		t.Fatalf("%v", err)
	}
	id := uint64(0)

	testObject := ht.MockHashable{id}

	hashableDao.Save(testObject)

	read, err := hashableDao.Get(id)

	if err != nil {

		t.Fatalf("%v", err)
	}
	assert.Equal(t, testObject, read)
}

func TestHashableDaoWithCollision(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	fileManager := fm.NewFileContainer(path, descriptor)

	daoId := uint64(0)

	hashableDao, err := ht.New[ht.MockHashable](fileManager, daoId)

	if err != nil {

		t.Fatalf("%v", err)
	}
	id := uint64(0)

	idB := uint64(10)

	idC := uint64(20)
	idD := uint64(30)
	idE := uint64(40)
	idF := uint64(50)

	testObject := ht.MockHashable{id}

	testObjectB := ht.MockHashable{idB}

	testObjectC := ht.MockHashable{idC}
	testObjectD := ht.MockHashable{idD}
	testObjectE := ht.MockHashable{idE}
	testObjectF := ht.MockHashable{idF}

	hashableDao.Save(testObject)

	hashableDao.Save(testObjectB)

	hashableDao.Save(testObjectC)
	hashableDao.Save(testObjectD)
	hashableDao.Save(testObjectE)
	hashableDao.Save(testObjectF)

	read, err := hashableDao.Get(id)

	readB, err := hashableDao.Get(idB)

	readC, err := hashableDao.Get(idC)
	readD, err := hashableDao.Get(idD)
	readE, err := hashableDao.Get(idE)
	readF, err := hashableDao.Get(idF)

	if err != nil {

		t.Fatalf("%v", err)
	}
	assert.Equal(t, testObject, read)

	assert.Equal(t, testObjectB, readB)

	assert.Equal(t, testObjectC, readC)
	assert.Equal(t, testObjectD, readD)
	assert.Equal(t, testObjectE, readE)
	assert.Equal(t, testObjectF, readF)
}

func TestHashableDaoWithCollisionAndDelete(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	fileManager := fm.NewFileContainer(path, descriptor)

	daoId := uint64(0)

	hashableDao, err := ht.New[ht.MockHashable](fileManager, daoId)

	if err != nil {

		t.Fatalf("%v", err)
	}
	id := uint64(0)

	idB := uint64(10)

	idC := uint64(20)

	idD := uint64(1)

	idE := uint64(11)

	testObject := ht.MockHashable{id}

	testObjectB := ht.MockHashable{idB}

	testObjectC := ht.MockHashable{idC}

	testObjectD := ht.MockHashable{idD}

	testObjectE := ht.MockHashable{idE}

	hashableDao.Save(testObject)

	hashableDao.Save(testObjectB)

	hashableDao.Save(testObjectC)

	hashableDao.Save(testObjectD)

	hashableDao.Save(testObjectE)

	hashableDao.Delete(id)

	read, err := hashableDao.Get(id)

	readB, err := hashableDao.Get(idB)

	readC, err := hashableDao.Get(idC)

	readE, err := hashableDao.Get(idE)

	assert.Nil(t, read)

	assert.Equal(t, testObjectB, readB)

	assert.Equal(t, testObjectC, readC)

	assert.Equal(t, testObjectE, readE)

	updatedTestObject := ht.MockHashable{id}

	updatedTestObjectD := ht.MockHashable{idD}

	err = hashableDao.Update(updatedTestObject)

	assert.Equal(t, ht.ObjectDoesNotExistToUpdate, err)

	err = hashableDao.Update(updatedTestObjectD)

	err = hashableDao.Delete(uint64(2))

	// this saves number 2 in the slot number one was deleted from but should not
	// because it already exists in a collision

	err = hashableDao.Save(testObjectB)

	readB, err = hashableDao.Get(idB)

	assert.Equal(t, testObjectB, readB)
}
