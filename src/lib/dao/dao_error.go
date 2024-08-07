package dao

type DaoErrorKind string

func (daoError DaoErrorKind) Error() string {

	return string(daoError)
}

const (
	ObjectAlreadyExists DaoErrorKind = DaoErrorKind("object already exists, cannot save new")

	ObjectDoesNotExistToUpdate DaoErrorKind = DaoErrorKind("object does not exist to update")

	ObjectDoesNotExistToDelete DaoErrorKind = DaoErrorKind("object does not exist to delete")
)
