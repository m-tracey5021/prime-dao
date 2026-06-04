package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/ht"
	"github.com/m-tracey5021/prime-dao/pkg/list"
)

func TestLinkedList(t *testing.T) {

	path := "dao_dir"

	defer fm.RemoveDir(path)

	descriptor := "obj_test"

	fileManager := fm.NewFileManager(path, descriptor)

	linkedList := list.NewList[ht.MockHashable](fileManager)

	testObjectC := ht.MockHashable{uuid.New(), 0}
	testObjectD := ht.MockHashable{uuid.New(), 0}
	testObjectE := ht.MockHashable{uuid.New(), 0}
	testObjectF := ht.MockHashable{uuid.New(), 0}

	testObjectG := ht.MockHashable{uuid.New(), 0}

	linkedList.Append(testObjectC)
	linkedList.Append(testObjectD)
	linkedList.Append(testObjectE)
	linkedList.Append(testObjectF)

	gotD, err := linkedList.Index(1)

	if err != nil {

		t.Fail()
	}
	fmt.Println(gotD)

	linkedList.Insert(testObjectG, 2) // C D G E F

	gotE, err := linkedList.Index(3)

	if err != nil {

		t.Fail()
	}
	fmt.Println(gotE)

	slice, err := linkedList.ToSlice()

	fmt.Println(slice)

	err = linkedList.Remove(1) // C G E F

	slice, err = linkedList.ToSlice()

	gotG, err := linkedList.Index(1)

	fmt.Println(gotG)

}
