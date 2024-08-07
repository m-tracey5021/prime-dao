package ht

type HashTableError string

func (err HashTableError) Error() string {

	return string(err)
}

const (
	ObjectAlreadyExists HashTableError = HashTableError("object already exists, cannot save new")

	ObjectDoesNotExistToUpdate HashTableError = HashTableError("object does not exist to update")

	ObjectDoesNotExistToDelete HashTableError = HashTableError("object does not exist to delete")

	CollisionAlreadyOccupied HashTableError = HashTableError("tried to save over an existing collision")
)
