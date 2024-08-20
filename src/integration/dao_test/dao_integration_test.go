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
	deletedSize, err := dao.ProcessDelete((*identifiableSavedB).Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", deletedSize)

	getC, err := dao.ProcessGet(identifiableSavedC.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", getC)

	fm.RemoveDir(path)
}
