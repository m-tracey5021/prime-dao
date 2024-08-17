package ht

type HashTableError string

func (err HashTableError) Error() string {

	return string(err)
}

const (
	ObjectAlreadyExists HashTableError = HashTableError("object already exists, cannot save new")

	ObjectDoesNotExist HashTableError = HashTableError("object does not exist")

	BucketOccupied HashTableError = HashTableError("bucket is already occupied")
)
