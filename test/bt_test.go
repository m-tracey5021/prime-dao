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

	testObjectDataA := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataB := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataC := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataD := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataE := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataF := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataG := MockIdentifiable{uuid.New(), []int{8}}

	testObjectA := ht.MockHashable{uuid.New(), 8}
	testObjectB := ht.MockHashable{uuid.New(), 9}
	testObjectC := ht.MockHashable{uuid.New(), 10}
	testObjectD := ht.MockHashable{uuid.New(), 11}
	testObjectE := ht.MockHashable{uuid.New(), 15}
	testObjectF := ht.MockHashable{uuid.New(), 20}
	testObjectG := ht.MockHashable{uuid.New(), 17}

	btree, err := bt.NewBTree[ht.MockHashable, MockIdentifiable](3, fileManager, func(key uuid.UUID) (*MockIdentifiable, error) {

		mapping := map[uuid.UUID]*MockIdentifiable{

			testObjectA.Id(): &testObjectDataA,
			testObjectB.Id(): &testObjectDataB,
			testObjectC.Id(): &testObjectDataC,
			testObjectD.Id(): &testObjectDataD,
			testObjectE.Id(): &testObjectDataE,
			testObjectF.Id(): &testObjectDataF,
			testObjectG.Id(): &testObjectDataG,
		}
		return mapping[key], nil
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

	result, _, _, _, err := btree.Search(testObjectE.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Println(result)
	if err := btree.Delete(testObjectD); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

}

func TestBTreeDeleteCase1(t *testing.T) {

	path := "dao_dir"

	defer fm.RemoveDir(path)

	descriptor := "obj_test"

	fileManager := fm.NewFileManager(path, descriptor)

	testObjectDataA := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataB := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataC := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataD := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataE := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataF := MockIdentifiable{uuid.New(), []int{8}}
	testObjectDataG := MockIdentifiable{uuid.New(), []int{8}}

	testObjectA := ht.MockHashable{uuid.New(), 8}
	testObjectB := ht.MockHashable{uuid.New(), 9}
	testObjectC := ht.MockHashable{uuid.New(), 10}
	testObjectD := ht.MockHashable{uuid.New(), 11}
	testObjectE := ht.MockHashable{uuid.New(), 15}
	testObjectF := ht.MockHashable{uuid.New(), 20}
	testObjectG := ht.MockHashable{uuid.New(), 17}

	btree, err := bt.NewBTree[ht.MockHashable, MockIdentifiable](3, fileManager, func(key uuid.UUID) (*MockIdentifiable, error) {

		mapping := map[uuid.UUID]*MockIdentifiable{

			testObjectA.Id(): &testObjectDataA,
			testObjectB.Id(): &testObjectDataB,
			testObjectC.Id(): &testObjectDataC,
			testObjectD.Id(): &testObjectDataD,
			testObjectE.Id(): &testObjectDataE,
			testObjectF.Id(): &testObjectDataF,
			testObjectG.Id(): &testObjectDataG,
		}
		return mapping[key], nil
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

	result, _, _, _, err := btree.Search(testObjectE.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Println(result)
	if err := btree.Delete(testObjectD); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

}
