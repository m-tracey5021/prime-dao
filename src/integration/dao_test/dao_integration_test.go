package dao_test

import (
	"testing"
	"transformer/src/lib/dao"

	"github.com/stretchr/testify/assert"
)

func TestHashableDao(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	hashableDao, err := dao.Wrapped[dao.MockHashable](path, descriptor)

	if err != nil {

		t.Fatalf("%v", err)
	}
	id := hashableDao.NewObjectId()

	testObject := dao.MockHashable{id}

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

	hashableDao, err := dao.Wrapped[dao.MockHashable](path, descriptor)

	if err != nil {

		t.Fatalf("%v", err)
	}
	id := uint64(0)

	idB := uint64(10)

	idC := uint64(20)

	testObject := dao.MockHashable{id}

	testObjectB := dao.MockHashable{idB}

	testObjectC := dao.MockHashable{idC}

	hashableDao.Save(testObject)

	hashableDao.Save(testObjectB)

	hashableDao.Save(testObjectC)

	read, err := hashableDao.Get(id)

	readB, err := hashableDao.Get(idB)

	readC, err := hashableDao.Get(idC)

	if err != nil {

		t.Fatalf("%v", err)
	}
	assert.Equal(t, testObject, read)

	assert.Equal(t, testObjectB, readB)

	assert.Equal(t, testObjectC, readC)
}
