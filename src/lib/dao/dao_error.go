package dao

type DaoErrorKind string

const (
	ObjectAlreadyExists DaoErrorKind = "object already exists, cannot save new"

	ObjectDoesNotExistToUpdate DaoErrorKind = "object does not exist to update"

	ObjectDoesNotExistToDelete DaoErrorKind = "object does not exist to delete"
)

type DaoError struct {
	err DaoErrorKind
}

func (daoError *DaoError) Error() string {

	return string(daoError.err)
}

func NewDaoError(err DaoErrorKind) *DaoError {

	return &DaoError{err}
}
