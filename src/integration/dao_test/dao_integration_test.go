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

	identifiableSaved, size, _ := dao.Save(identifiable)

	identifiableSavedB, sizeB, _ := dao.Save(identifiableB)

	identifiableSavedC, sizeC, _ := dao.Save(identifiableC)

	identifiableSavedD, sizeD, _ := dao.Save(identifiableD)

	indexes, err := dao.GetAllIndexes()

	fmt.Printf("%v", indexes)

	fmt.Printf("%v, size %v", identifiableSaved, size)
	fmt.Printf("%v, size %v", identifiableSavedB, sizeB)
	fmt.Printf("%v, size %v", identifiableSavedC, sizeC)
	fmt.Printf("%v, size %v", identifiableSavedD, sizeD)

	getB, err := dao.Get(identifiableSavedB.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", getB)

	identifiableSavedB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	updatedSize, err := dao.Update(*identifiableSavedB)

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v, size %v", identifiableSavedB, updatedSize)

	getB, err = dao.Get(identifiableSavedB.Id())

	fmt.Printf("%v", getB)

	if err != nil {

		t.Fail()
	}
	_, err = dao.Delete((*identifiableSaved).Id())

	deletedSizeB, err := dao.Delete((*identifiableSavedB).Id())

	if err != nil {

		t.Fail()
	}
	deletedSizeC, err := dao.Delete((*identifiableSavedC).Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", deletedSizeB)
	fmt.Printf("%v", deletedSizeC)

	getC, err := dao.Get(identifiableSavedC.Id())

	// should be err here
	fmt.Printf("%v", getC)

	getD, err := dao.Get(identifiableSavedD.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", getD)

	identifiableE := MockIdentifiable{Data: []int{0, 2, 9}}

	identifiableSavedE, _, _ := dao.Save(identifiableE)

	fmt.Printf("%v", identifiableSavedE)

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

	saveTransaction := dao.NewTransaction()

	_ = saveTransaction.Save(identifiable)
	requestB := saveTransaction.Save(identifiableB)
	requestC := saveTransaction.Save(identifiableC)
	requestD := saveTransaction.Save(identifiableD)

	results := dao.ExecuteTransaction(saveTransaction)

	identifiableB = *results[requestB].Object()

	identifiableC = *results[requestC].Object()

	identifiableD = *results[requestD].Object()

	fmt.Printf("%v", identifiableB)

	transaction := dao.NewTransaction()

	getRequest := transaction.Get(identifiableB.Id())

	identifiableB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	requestE := transaction.WithDependency(getRequest).Update(identifiableB)

	requestF := transaction.WithDependency(requestE).Get(identifiableB.Id())

	requestG := transaction.WithDependency(requestF).Delete(identifiable.Id())

	requestGb := transaction.WithDependency(requestF).Delete(identifiableB.Id())

	if err != nil {

		t.Fail()
	}
	requestH := transaction.Delete(identifiableC.Id())

	requestI := transaction.WithDependency(requestH).Get(identifiableC.Id())

	requestJ := transaction.Get(identifiableD.Id())

	identifiableE := MockIdentifiable{Data: []int{0, 2, 9}}

	requestK := transaction.Save(identifiableE)

	results = dao.ExecuteTransaction(transaction)

	resultE := results[requestE]
	resultF := results[requestF]
	resultG := results[requestG]
	resultGb := results[requestGb]
	resultH := results[requestH]
	resultI := results[requestI]
	resultJ := results[requestJ]
	resultK := results[requestK]

	fmt.Printf("%v", resultE)
	fmt.Printf("%v", resultF)
	fmt.Printf("%v", resultG)
	fmt.Printf("%v", resultGb)
	fmt.Printf("%v", resultH)
	fmt.Printf("%v", resultI)
	fmt.Printf("%v", resultJ)
	fmt.Printf("%v", resultK)

	fm.RemoveDir(path)
}
