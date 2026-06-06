package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/dao"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type MockIdentifiable struct {
	MockId uuid.UUID

	Data []int
}

func (identifiable MockIdentifiable) Id() uuid.UUID {

	return identifiable.MockId
}

func (identifiable MockIdentifiable) SetId(id uuid.UUID) {

	identifiable.MockId = id
}

func (identifiable MockIdentifiable) Descriptor() string {

	return "mock"
}

func (identifiable MockIdentifiable) SortKeyValue() int {

	return 1
}

func (identifiable MockIdentifiable) Compare(other schema.Orderable[int]) schema.Order {

	return schema.Equal
}

func TestDao(t *testing.T) {

	path := "dao_dir"

	// defer fm.RemoveDir(path)

	descriptor := "obj_test"

	dao, err := dao.From[MockIdentifiable, int](fm.NewFileManager(path, descriptor))

	if err != nil {

		t.Fail()
	}
	identifiable := MockIdentifiable{Data: []int{1, 2}}
	identifiableB := MockIdentifiable{Data: []int{2, 3, 4}}
	identifiableC := MockIdentifiable{Data: []int{1, 4}}
	identifiableD := MockIdentifiable{Data: []int{0, 2}}

	dao.Save(identifiable)

	dao.Save(identifiableB)

	dao.Save(identifiableC)

	dao.Save(identifiableD)

	getA, err := dao.Get(identifiable.Id())
	getB, err := dao.Get(identifiableB.Id())
	getC, err := dao.Get(identifiableC.Id())
	getD, err := dao.Get(identifiableD.Id())

	getA, err = dao.Get(identifiable.Id())
	getC, err = dao.Get(identifiableC.Id())

	if err != nil {

		t.Fail()
	}

	fmt.Printf("%v", getA)
	fmt.Printf("%v", getB)
	fmt.Printf("%v", getC)
	fmt.Printf("%v", getD)

	identifiableB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	updatedSize, err := dao.Update(identifiableB)

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v, size %v", identifiableB, updatedSize)

	getB, err = dao.Get(identifiableB.Id())

	fmt.Printf("%v", getB)

	if err != nil {

		t.Fail()
	}
	_, err = dao.Delete(identifiable.Id())

	deletedSizeB, err := dao.Delete(identifiableB.Id())

	if err != nil {

		t.Fail()
	}
	deletedSizeC, err := dao.Delete(identifiableC.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", deletedSizeB)
	fmt.Printf("%v", deletedSizeC)

	getC, err = dao.Get(identifiableC.Id())

	// should be err here
	fmt.Printf("%v", getC)

	getD, err = dao.Get(identifiableD.Id())

	if err != nil {

		t.Fail()
	}
	fmt.Printf("%v", getD)

	identifiableE := MockIdentifiable{Data: []int{0, 2, 9}}

	dao.Save(identifiableE)

	fmt.Printf("%v", identifiableE)
}

func TestDaoPersists(t *testing.T) {

	path := "dao_dir"

	defer fm.RemoveDir(path)

	descriptor := "obj_test"

	dao, err := dao.From[MockIdentifiable](fm.NewFileManager(path, descriptor))

	if err != nil {

		t.Fail()
	}
	all, err := dao.GetAll()

	fmt.Println(all)
}

func TestDaoAsyncSimple(t *testing.T) {

	path := "dao_dir"

	// defer fm.RemoveDir(path)

	descriptor := "obj_test"

	dao, err := dao.From[MockIdentifiable](fm.NewFileManager(path, descriptor))

	if err != nil {

		t.Fail()
	}
	identifiable := MockIdentifiable{Data: []int{1, 2}}
	identifiableB := MockIdentifiable{Data: []int{2, 3, 4}}
	identifiableC := MockIdentifiable{Data: []int{1, 4}}
	identifiableD := MockIdentifiable{Data: []int{0, 2}}

	saveTransaction := dao.NewTransaction()

	saveTransaction.Save(identifiable)
	saveTransaction.Save(identifiableB)
	saveTransaction.Save(identifiableC)
	saveTransaction.Save(identifiableD)

	results, _ := dao.ExecuteTransaction(saveTransaction)

	// fmt.Printf("%v\n", results)

	transaction := dao.NewTransaction()

	getRequest := transaction.Get(identifiableB.Id())

	identifiableB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	requestE := transaction.WithDependency(getRequest).Update(identifiableB)

	requestF := transaction.WithDependency(requestE).Get(identifiableB.Id())

	results, _ = dao.ExecuteTransaction(transaction)

	resultE := results[requestE]
	resultF := results[requestF]

	y, err := dao.Get(identifiableB.Id())

	fmt.Printf("%v\n", y)

	fmt.Printf("%v\n", resultE)
	fmt.Printf("%v\n", resultF)
}

func TestDaoAsync(t *testing.T) {

	path := "dao_dir"

	defer fm.RemoveDir(path)

	descriptor := "obj_test"

	dao, err := dao.From[MockIdentifiable](fm.NewFileManager(path, descriptor))

	if err != nil {

		t.Fail()
	}
	identifiable := MockIdentifiable{Data: []int{1, 2}}
	identifiableB := MockIdentifiable{Data: []int{2, 3, 4}}
	identifiableC := MockIdentifiable{Data: []int{1, 4}}
	identifiableD := MockIdentifiable{Data: []int{0, 2}}

	saveTransaction := dao.NewTransaction()

	saveTransaction.Save(identifiable)
	saveTransaction.Save(identifiableB)
	saveTransaction.Save(identifiableC)
	saveTransaction.Save(identifiableD)

	results, _ := dao.ExecuteTransaction(saveTransaction)

	fmt.Printf("%v", results)

	transaction := dao.NewTransaction()

	getRequest := transaction.Get(identifiableB.Id())

	identifiableB.Data = []int{2, 3, 4, 5, 6, 7, 8, 1, 2, 3, 4, 5, 6, 7, 8}

	requestE := transaction.WithDependency(getRequest).Update(identifiableB)

	requestF := transaction.WithDependency(requestE).Get(identifiableB.Id())

	requestG := transaction.Delete(identifiable.Id())

	requestGb := transaction.WithDependency(requestF).Delete(identifiableB.Id())

	if err != nil {

		t.Fail()
	}
	requestH := transaction.Delete(identifiableC.Id())

	requestI := transaction.WithDependency(requestH).Get(identifiableC.Id())

	requestJ := transaction.Get(identifiableD.Id())

	identifiableE := MockIdentifiable{Data: []int{0, 2, 9}}

	requestK := transaction.Save(identifiableE)

	results, _ = dao.ExecuteTransaction(transaction)

	resultE := results[requestE]
	resultF := results[requestF]
	resultG := results[requestG]
	resultGb := results[requestGb]
	resultH := results[requestH]
	resultI := results[requestI]
	resultJ := results[requestJ]
	resultK := results[requestK]

	x, err := dao.Delete(identifiableB.Id())

	y, err := dao.Get(identifiableB.Id())

	fmt.Printf("%v", x)
	fmt.Printf("%v", y)

	fmt.Printf("%v", resultE)
	fmt.Printf("%v", resultF)
	fmt.Printf("%v", resultG)
	fmt.Printf("%v", resultGb)
	fmt.Printf("%v", resultH)
	fmt.Printf("%v", resultI)
	fmt.Printf("%v", resultJ)
	fmt.Printf("%v", resultK)
}
