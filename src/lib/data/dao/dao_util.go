package dao

import "prime-dao/src/lib/data/schema"

func GetSome[T schema.DescribedIdentifiable](dao TSFDao[T], ids ...uint64) ([]T, error) {

	transaction := dao.NewTransaction()

	for _, id := range ids {

		transaction.Get(id)
	}
	result := dao.ExecuteTransaction(transaction)

	objects := make([]T, 0)

	for _, nthResult := range result {

		err := nthResult.Error()

		if err != nil {

			return objects, err
		}
		objects = append(objects, *nthResult.Object())
	}
	return objects, nil
}

func GetAll[T schema.DescribedIdentifiable](dao TSFDao[T]) ([]T, error) {

	ids := dao.AllObjectIds()

	return GetSome(dao, ids...)
}

func GetByName[T schema.Named](dao TSFDao[T], name string) (*T, error) {

	objects, err := GetAll(dao)

	if err != nil {

		return nil, err
	}
	for _, object := range objects {

		if object.Name() == name {

			return &object, nil
		}
	}
	return nil, nil
}
