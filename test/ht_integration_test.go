package test

// import (
// 	"io/ioutil"
// 	"os"
// 	"path/filepath"
// 	"testing"
// 	"time"
// 	"github.com/m-tracey5021/prime-dao/src/lib/dao/ht"

// 	"github.com/stretchr/testify/assert"
// )

// func TestHashableDao(t *testing.T) {

// 	path := "dao_dir"

// 	descriptor := "obj_test"

// 	// fileManager := fm.NewFileContainer(path, descriptor)

// 	daoId := uint64(0)

// 	hashableDao := ht.New[ht.MockHashable](path, descriptor, daoId)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	id := uint64(0)

// 	testObject := ht.MockHashable{id}

// 	hashableDao.Save(testObject)

// 	read, err := hashableDao.Get(id)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	assert.Equal(t, testObject, read)
// }

// func TestHashableDaoWithCollision(t *testing.T) {

// 	path := "dao_dir"

// 	descriptor := "obj_test"

// 	// fileManager := fm.NewFileContainer(path, descriptor)

// 	daoId := uint64(0)

// 	hashableDao, err := ht.New[ht.MockHashable](path, descriptor, daoId)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	id := uint64(0)

// 	idB := uint64(10)

// 	idC := uint64(20)
// 	idD := uint64(30)
// 	idE := uint64(40)
// 	idF := uint64(50)

// 	testObject := ht.MockHashable{id}

// 	testObjectB := ht.MockHashable{idB}

// testObjectC := ht.MockHashable{idC}
// testObjectD := ht.MockHashable{idD}
// testObjectE := ht.MockHashable{idE}
// testObjectF := ht.MockHashable{idF}

// 	hashableDao.Save(testObject)

// 	hashableDao.Save(testObjectB)

// 	hashableDao.Save(testObjectC)
// 	hashableDao.Save(testObjectD)
// 	hashableDao.Save(testObjectE)
// 	hashableDao.Save(testObjectF)

// 	read, err := hashableDao.Get(id)

// 	readB, err := hashableDao.Get(idB)

// 	readC, err := hashableDao.Get(idC)
// 	readD, err := hashableDao.Get(idD)
// 	readE, err := hashableDao.Get(idE)
// 	readF, err := hashableDao.Get(idF)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	assert.Equal(t, testObject, read)

// 	assert.Equal(t, testObjectB, readB)

// 	assert.Equal(t, testObjectC, readC)
// 	assert.Equal(t, testObjectD, readD)
// 	assert.Equal(t, testObjectE, readE)
// 	assert.Equal(t, testObjectF, readF)
// }

// func TestHashableDaoWithCollisionAndDelete(t *testing.T) {

// 	path := "dao_dir"

// 	descriptor := "obj_test"

// 	// fileManager := fm.NewFileContainer(path, descriptor)

// 	daoId := uint64(0)

// 	hashableDao, err := ht.New[ht.MockHashable](path, descriptor, daoId)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	id := uint64(0)

// 	idB := uint64(10)

// 	idC := uint64(20)

// 	idD := uint64(1)

// 	idE := uint64(11)

// 	testObject := ht.MockHashable{id}

// 	testObjectB := ht.MockHashable{idB}

// 	testObjectC := ht.MockHashable{idC}

// 	testObjectD := ht.MockHashable{idD}

// 	testObjectE := ht.MockHashable{idE}

// 	hashableDao.Save(testObject)

// 	hashableDao.Save(testObjectB)

// 	hashableDao.Save(testObjectC)

// 	hashableDao.Save(testObjectD)

// 	hashableDao.Save(testObjectE)

// 	hashableDao.Delete(id)

// 	read, err := hashableDao.Get(id)

// 	readB, err := hashableDao.Get(idB)

// 	readC, err := hashableDao.Get(idC)

// 	readE, err := hashableDao.Get(idE)

// 	assert.Nil(t, read)

// 	assert.Equal(t, testObjectB, readB)

// 	assert.Equal(t, testObjectC, readC)

// 	assert.Equal(t, testObjectE, readE)

// 	updatedTestObject := ht.MockHashable{id}

// 	updatedTestObjectD := ht.MockHashable{idD}

// 	err = hashableDao.Update(updatedTestObject)

// 	assert.Equal(t, ht.ObjectDoesNotExist, err)

// 	err = hashableDao.Update(updatedTestObjectD)

// 	err = hashableDao.Delete(uint64(2))

// 	// this saves number 2 in the slot number one was deleted from but should not
// 	// because it already exists in a collision

// 	err = hashableDao.Save(testObjectB)

// 	readB, err = hashableDao.Get(idB)

// 	assert.Equal(t, testObjectB, readB)

// 	removeAllFiles(path)
// }

// func TestHashableDaoWithManySaves(t *testing.T) {

// 	path := "dao_dir"

// 	descriptor := "obj_test"

// 	// fileManager := fm.NewFileContainer(path, descriptor)

// 	daoId := uint64(0)

// 	hashTable, err := ht.New[ht.MockHashable](path, descriptor, daoId)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	ids := make([]uint64, 0)

// 	objects := make([]ht.MockHashable, 0)

// 	for number := range 10000 {

// 		id := uint64(number)

// 		ids = append(ids, uint64(id))

// 		testObject := ht.MockHashable{id}

// 		objects = append(objects, testObject)

// 	}
// 	start := time.Now()

// 	for _, obj := range objects {

// 		hashTable.Save(obj)
// 	}
// 	for _, id := range ids {

// 		hashTable.Get(id)
// 	}

// 	duration := time.Since(start)
// 	t.Logf("Get in a loop took %s to run", duration)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}

// 	removeAllFiles(path)

// }

// func TestHashableDaoWithManySavesConcurrent(t *testing.T) {

// 	path := "dao_dir"

// 	descriptor := "obj_test"

// 	// fileManager := fm.NewFileContainer(path, descriptor)

// 	daoId := uint64(0)

// 	hashTable, err := ht.New[ht.MockHashable](path, descriptor, daoId)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	ids := make([]uint64, 0)

// 	objects := make([]ht.MockHashable, 0)

// 	for number := range 10000 {

// 		id := uint64(number)

// 		ids = append(ids, id)

// 		objects = append(objects, ht.MockHashable{id})
// 	}
// 	start := time.Now()

// 	hashTable.SaveAsync(objects...)

// 	hashTable.GetAsync(ids...)

// 	_ = hashTable.GetResults()

// 	duration := time.Since(start)

// 	// t.Logf("Get and Save Concurrent results: %s", results)

// 	t.Logf("Get and Save Concurrent took %s to run", duration)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	// assert.NotNil(t, reads)

// 	removeAllFiles(path)

// }

// func TestHashableDaoWithAFewSavesConcurrent(t *testing.T) {

// 	path := "dao_dir"

// 	descriptor := "obj_test"

// 	// fileManager := fm.NewFileContainer(path, descriptor)

// 	daoId := uint64(0)

// 	hashTable, err := ht.New[ht.MockHashable](path, descriptor, daoId)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	ids := make([]uint64, 0)

// 	objects := make([]ht.MockHashable, 0)

// 	for number := range 20 {

// 		id := uint64(number)

// 		ids = append(ids, id)

// 		objects = append(objects, ht.MockHashable{id})
// 	}
// 	start := time.Now()

// 	hashTable.SaveAsync(objects...)

// 	hashTable.GetAsync(ids...)

// 	results := hashTable.GetResults()

// 	duration := time.Since(start)

// 	t.Logf("Get and Save Concurrent results: %s", results)

// 	t.Logf("Get and Save Concurrent took %s to run", duration)

// 	if err != nil {

// 		t.Fatalf("%v", err)
// 	}
// 	// assert.NotNil(t, reads)

// 	removeAllFiles(path)

// }

// func removeAllFiles(dir string) error {
// 	// Get a list of all files in the directory
// 	files, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		return err
// 	}

// 	// Loop through all files and remove them
// 	for _, file := range files {
// 		filePath := filepath.Join(dir, file.Name())
// 		err := os.Remove(filePath)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }
