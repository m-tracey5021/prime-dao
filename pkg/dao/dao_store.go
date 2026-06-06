package dao

import (
	"errors"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type ComponentStore[T schema.Orderable[U], U constraints.Ordered] struct {
	objectsOnDisk map[uuid.UUID]T

	objectsInMemory map[uuid.UUID]T

	saveRequests map[uuid.UUID]uint64

	updateRequests map[uuid.UUID]uint64

	objectDao *Dao[T, U]

	transaction DaoTransaction[T, U]
}

func NewComponentStore[T schema.Orderable[U], U constraints.Ordered](dao *Dao[T, U]) ComponentStore[T, U] {

	return ComponentStore[T, U]{
		objectsOnDisk: map[uuid.UUID]T{},

		objectsInMemory: map[uuid.UUID]T{},

		saveRequests: map[uuid.UUID]uint64{},

		updateRequests: map[uuid.UUID]uint64{},

		objectDao: dao,

		transaction: dao.NewTransaction(),
	}
}

func (manager *ComponentStore[T, U]) Objects() map[uuid.UUID]T {

	return manager.objectsOnDisk
}

func (manager *ComponentStore[T, U]) AddObject(object T) {

	object.SetId(uuid.New()) // set a tmp id, this will be overwritten when saved

	requestId := manager.transaction.Save(object)

	manager.objectsInMemory[object.Id()] = object

	manager.saveRequests[object.Id()] = requestId
}

func (manager *ComponentStore[T, U]) GetCachedObject(id uuid.UUID) (*T, bool) {

	object, ok := manager.objectsOnDisk[id]

	if ok {

		return &object, true
	}
	objectInMem, ok := manager.objectsInMemory[id]

	if ok {

		return &objectInMem, false
	}
	return nil, false
}

func (manager *ComponentStore[T, U]) ReadObject(id uuid.UUID) (*T, error) {

	cached, _ := manager.GetCachedObject(id)

	if cached != nil {

		return cached, nil

	} else {

		transaction := manager.objectDao.NewTransaction() // dont even need to use transaction here??

		request := transaction.Get(id)

		executionResult, err := manager.objectDao.ExecuteTransaction(transaction)

		if err != nil {

			return nil, err
		}
		result := executionResult[request]

		return result.Object(), result.Error()
	}
}

func (manager *ComponentStore[T, U]) ReadObjects(ids ...uuid.UUID) ([]T, error) {

	objects := []T{}

	notFound := []uuid.UUID{}

	for _, id := range ids {

		cached, _ := manager.GetCachedObject(id)

		if cached != nil {

			objects = append(objects, *cached)

		} else {

			notFound = append(notFound, id)
		}
	}
	if len(notFound) == 0 {

		return objects, nil
	}
	transaction := manager.objectDao.NewTransaction()

	for _, id := range notFound {

		transaction.Get(id)
	}
	result, err := manager.objectDao.ExecuteTransaction(transaction)

	if err != nil {

		return []T{}, err
	}
	for _, nthResult := range result {

		err := nthResult.Error()

		if err != nil {

			return objects, err
		}
		object := *nthResult.Object()

		objects = append(objects, object)

		manager.objectsOnDisk[object.Id()] = object
	}
	return objects, nil
}

// func (manager *ComponentStore[T, U]) ReadAll() ([]T, error) {

// 	ids := manager.objectDao.AllObjectIds()

// 	return manager.ReadObjects(ids...)
// }

// func (manager *ComponentStore[T, U]) ReadObjectIntoMemByCondition(condition func(T) bool) (*T, error) {

// 	for _, object := range manager.Objects() {

// 		if condition(object) {

// 			return &object, nil
// 		}
// 	}
// 	objects, err := manager.ReadAll()

// 	if err != nil {

// 		return nil, err
// 	}
// 	for _, object := range objects {

// 		if condition(object) {

// 			return &object, nil
// 		}
// 	}
// 	return nil, nil
// }

func (manager *ComponentStore[T, U]) UpdateObject(object T) {

	cached, onDisk := manager.GetCachedObject(object.Id())

	if cached != nil {

		if onDisk {

			updateRequestId, existingUpdateRequest := manager.updateRequests[object.Id()]

			if existingUpdateRequest {

				manager.transaction.RemoveRequest(updateRequestId)

				updatedUpdateRequest := manager.transaction.Update(object)

				manager.updateRequests[object.Id()] = updatedUpdateRequest

			} else {

				manager.transaction.Update(object)
			}
			manager.objectsOnDisk[object.Id()] = object

		} else {

			saveRequestId, existingSaveRequest := manager.saveRequests[object.Id()]

			if existingSaveRequest {

				manager.transaction.RemoveRequest(saveRequestId)

				updatedSaveRequest := manager.transaction.Save(object)

				manager.saveRequests[object.Id()] = updatedSaveRequest

			} else {

				manager.transaction.Save(object)
			}
			manager.objectsInMemory[object.Id()] = object
		}

	} else {

		manager.transaction.Update(object)
	}
}

func (manager *ComponentStore[T, U]) DeleteObject(object T) {

	cached, onDisk := manager.GetCachedObject(object.Id())

	if cached != nil {

		if onDisk {

			updateRequestId, existingUpdateRequest := manager.updateRequests[object.Id()]

			if existingUpdateRequest {

				manager.transaction.RemoveRequest(updateRequestId)

				delete(manager.updateRequests, object.Id())
			}
			delete(manager.objectsOnDisk, object.Id())

		} else {

			saveRequestId, existingSaveRequest := manager.saveRequests[object.Id()]

			if existingSaveRequest {

				manager.transaction.RemoveRequest(saveRequestId)

				delete(manager.saveRequests, object.Id())
			}
			delete(manager.objectsInMemory, object.Id())
		}

	} else {

		manager.transaction.Delete(object.Id())
	}
}

func (manager *ComponentStore[T, U]) DeletePossibleUpdateRequest(objectId uuid.UUID) bool {

	updateRequestId, existingUpdate := manager.updateRequests[objectId]

	if existingUpdate {

		manager.transaction.RemoveRequest(updateRequestId)

		delete(manager.updateRequests, objectId)

		return true
	}
	return false
}

func (manager *ComponentStore[T, U]) DeletePossibleSaveRequest(objectId uuid.UUID) bool {

	saveRequestId, existingSave := manager.saveRequests[objectId]

	if existingSave {

		manager.transaction.RemoveRequest(saveRequestId)

		delete(manager.saveRequests, objectId)

		return true
	}
	return false
}

func (manager *ComponentStore[T, U]) Commit() error {

	result, err := manager.objectDao.ExecuteTransaction(manager.transaction)

	if err != nil {

		return err
	}
	var joinedErr error

	for _, result := range result {

		if result.Error() != nil {

			joinedErr = errors.Join(joinedErr, result.Error())
		}
	}
	return joinedErr
}
