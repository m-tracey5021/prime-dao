package test

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/bt"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/ht"
)

func UUIDFromInt(id uint64) uuid.UUID {
	var b [16]byte
	binary.BigEndian.PutUint64(b[8:], id) // Store in the last 8 bytes
	return uuid.UUID(b)
}

func TestBTree(t *testing.T) {

	path := "dao_dir"

	defer fm.RemoveDir(path)

	descriptor := "obj_test"

	fileManager := fm.NewFileManager(path, descriptor)

	btree, err := bt.NewBTree[*ht.MockHashable](3, fileManager, func(a, b *ht.MockHashable) int { return 0 })

	if err != nil {

		t.Fail()
	}

	testObjectA := ht.MockHashable{uuid.New(), 8}
	testObjectB := ht.MockHashable{uuid.New(), 9}
	testObjectC := ht.MockHashable{uuid.New(), 10}
	testObjectD := ht.MockHashable{uuid.New(), 11}
	testObjectE := ht.MockHashable{uuid.New(), 15}
	testObjectF := ht.MockHashable{uuid.New(), 20}
	testObjectG := ht.MockHashable{uuid.New(), 17}

	if err := btree.Insert(&testObjectA); err != nil {

		t.Fail()
	}
	str := btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectB); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectC); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectD); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectE); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectF); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectG); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

	result, _, _, _, err := btree.Search(&testObjectE)

	if err != nil {

		t.Fail()
	}
	fmt.Println(result)
	if err := btree.Delete(&testObjectD); err != nil {

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

	btree, err := bt.NewBTree[*ht.MockHashable](3, fileManager, func(a, b *ht.MockHashable) int { return 0 })

	if err != nil {

		t.Fail()
	}

	testObjectA := ht.MockHashable{uuid.New(), 8}
	testObjectB := ht.MockHashable{uuid.New(), 9}
	testObjectC := ht.MockHashable{uuid.New(), 10}
	testObjectD := ht.MockHashable{uuid.New(), 11}
	testObjectE := ht.MockHashable{uuid.New(), 15}
	testObjectF := ht.MockHashable{uuid.New(), 20}
	testObjectG := ht.MockHashable{uuid.New(), 17}

	if err := btree.Insert(&testObjectA); err != nil {

		t.Fail()
	}
	str := btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectB); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectC); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectD); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectE); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectF); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)
	if err := btree.Insert(&testObjectG); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

	result, _, _, _, err := btree.Search(&testObjectE)

	if err != nil {

		t.Fail()
	}
	fmt.Println(result)
	if err := btree.Delete(&testObjectD); err != nil {

		t.Fail()
	}
	str = btree.ToString()

	fmt.Println(str)

}
