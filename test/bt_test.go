package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/bt"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/ht"
)

func TestBTree(t *testing.T) {

	path := "dao_dir"

	defer fm.RemoveDir(path)

	descriptor := "obj_test"

	fileManager := fm.NewFileManager(path, descriptor)

	testObjectA := ht.MockHashable{uuid.New(), 8}
	testObjectB := ht.MockHashable{uuid.New(), 9}
	testObjectC := ht.MockHashable{uuid.New(), 10}
	testObjectD := ht.MockHashable{uuid.New(), 11}
	testObjectE := ht.MockHashable{uuid.New(), 15}
	testObjectF := ht.MockHashable{uuid.New(), 20}
	testObjectG := ht.MockHashable{uuid.New(), 17}

	btree, err := bt.NewBTree[ht.MockHashable, int](3, fileManager, func(key uuid.UUID) (*int, error) {

		mapping := map[uuid.UUID]int{

			testObjectA.Id(): 1,
			testObjectB.Id(): 2,
			testObjectC.Id(): 3,
			testObjectD.Id(): 4,
			testObjectE.Id(): 5,
			testObjectF.Id(): 6,
			testObjectG.Id(): 7,
		}
		result := mapping[key]

		return &result, nil
	})

	if err != nil {

		t.Fail()
	}

	if err := btree.Insert(testObjectA); err != nil {

		t.Fail()
	}
	str := btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(testObjectB); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(testObjectC); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(testObjectD); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(testObjectE); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(testObjectF); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(testObjectG); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

	result, _, _, _, err := btree.Search(5)

	if err != nil {

		t.Fail()
	}
	fmt.Println(result)
	if err := btree.Delete(4); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

}
