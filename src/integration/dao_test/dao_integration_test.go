package dao_test

import (
	"fmt"
	"testing"
	"transformer/src/lib/data/dao"
	"transformer/src/lib/data/fm"
	"transformer/src/lib/data/schema"
)

type MockIdentifiable struct {
	MockId uint64

	Data []int
}

func (identifiable MockIdentifiable) Id() uint64 {

	return identifiable.MockId
}

func (identifiable MockIdentifiable) SetId(id uint64) schema.Identifiable {

	return MockIdentifiable{id, identifiable.Data}
}

func TestDao(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	dao, err := dao.New[MockIdentifiable](path, descriptor, 0)

	if err != nil {

		t.Fail()
	}
	identifiable := MockIdentifiable{Data: []int{1, 2}}
	identifiableB := MockIdentifiable{Data: []int{2, 3, 4}}
	identifiableC := MockIdentifiable{Data: []int{1, 4}}
	identifiableD := MockIdentifiable{Data: []int{0, 2}}

	identifiableSaved, size, _ := dao.ProcessSave(identifiable)

	identifiableSavedB, sizeB, _ := dao.ProcessSave(identifiableB)

	identifiableSavedC, sizeC, _ := dao.ProcessSave(identifiableC)

	identifiableSavedD, sizeD, _ := dao.ProcessSave(identifiableD)

	fmt.Printf("%v, size %v", identifiableSaved, size)
	fmt.Printf("%v, size %v", identifiableSavedB, sizeB)
	fmt.Printf("%v, size %v", identifiableSavedC, sizeC)
	fmt.Printf("%v, size %v", identifiableSavedD, sizeD)

	getB, err := dao.ProcessGet(identifiableSavedB.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", getB)

	identifiableSavedB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	updatedSize, err := dao.ProcessUpdate(*identifiableSavedB)

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v, size %v", identifiableSavedB, updatedSize)

	getB, err = dao.ProcessGet(identifiableSavedB.Id())

	fmt.Printf("%v", getB)

	if err != nil {

		t.Fail()
	}
	deletedSizeB, err := dao.ProcessDelete((*identifiableSavedB).Id())

	if err != nil {

		t.Fail()
	}
	deletedSizeC, err := dao.ProcessDelete((*identifiableSavedC).Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", deletedSizeB)
	fmt.Printf("%v", deletedSizeC)

	getC, err := dao.ProcessGet(identifiableSavedC.Id())

	// should be err here
	fmt.Printf("%v", getC)

	getD, err := dao.ProcessGet(identifiableSavedD.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", getD)

	identifiableE := MockIdentifiable{Data: []int{0, 2, 9}}

	identifiableSavedE, _, _ := dao.ProcessSave(identifiableE)

	fmt.Printf("%v", identifiableSavedE)

	fm.RemoveDir(path)
}

func TestDaoContextAsync(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	dao, err := dao.New[MockIdentifiable](path, descriptor, 0)

	if err != nil {

		t.Fail()
	}
	identifiable := MockIdentifiable{Data: []int{1, 2}}
	identifiableB := MockIdentifiable{Data: []int{2, 3, 4}}
	identifiableC := MockIdentifiable{Data: []int{1, 4}}
	identifiableD := MockIdentifiable{Data: []int{0, 2}}

	_ = dao.QueueSaveRequest(identifiable)
	requestB := dao.QueueSaveRequest(identifiableB)
	requestC := dao.QueueSaveRequest(identifiableC)
	_ = dao.QueueSaveRequest(identifiableD)

	dao.QueueGetRequest(0, requestB)
	dao.QueueGetRequest(1, requestC)

	results := dao.ExecuteReq()

	fmt.Printf("%v", results)

	fm.RemoveDir(path)
}

func TestDaoAsync(t *testing.T) {

	path := "dao_dir"

	descriptor := "obj_test"

	dao, err := dao.New[MockIdentifiable](path, descriptor, 0)

	if err != nil {

		t.Fail()
	}
	identifiable := MockIdentifiable{Data: []int{1, 2}}
	identifiableB := MockIdentifiable{Data: []int{2, 3, 4}}
	identifiableC := MockIdentifiable{Data: []int{1, 4}}
	identifiableD := MockIdentifiable{Data: []int{0, 2}}

	_ = dao.QueueSaveRequest(identifiable)
	requestB := dao.QueueSaveRequest(identifiableB)
	requestC := dao.QueueSaveRequest(identifiableC)
	requestD := dao.QueueSaveRequest(identifiableD)

	results := dao.Execute()

	identifiableB = *results[requestB].Object()

	identifiableC = *results[requestC].Object()

	identifiableD = *results[requestD].Object()

	fmt.Printf("%v", identifiableB)

	dao.QueueGetRequest(identifiableB.Id(), requestB)

	identifiableB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	requestE := dao.QueueUpdateRequest(identifiableB)

	requestF := dao.QueueGetRequest(identifiableB.Id())

	requestG := dao.QueueDeleteRequest(identifiableB.Id())

	if err != nil {

		t.Fail()
	}
	requestH := dao.QueueDeleteRequest(identifiableC.Id())

	requestI := dao.QueueGetRequest(identifiableC.Id())

	requestJ := dao.QueueGetRequest(identifiableD.Id())

	identifiableE := MockIdentifiable{Data: []int{0, 2, 9}}

	requestK := dao.QueueSaveRequest(identifiableE)

	results = dao.Execute()

	resultE := results[requestE]
	resultF := results[requestF]
	resultG := results[requestG]
	resultH := results[requestH]
	resultI := results[requestI]
	resultJ := results[requestJ]
	resultK := results[requestK]

	fmt.Printf("%v", resultE)
	fmt.Printf("%v", resultF)
	fmt.Printf("%v", resultG)
	fmt.Printf("%v", resultH)
	fmt.Printf("%v", resultI)
	fmt.Printf("%v", resultJ)
	fmt.Printf("%v", resultK)

	fm.RemoveDir(path)
}
