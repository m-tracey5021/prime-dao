package ht

type HashTableError string

func (err HashTableError) Error() string {

	return string(err)
}

const (
	ObjectAlreadyExists HashTableError = HashTableError("object already exists, cannot save new")

	ObjectDoesNotExist HashTableError = HashTableError("object does not exist")
	ObjectNotFound     HashTableError = HashTableError("object does not exist")

	// If this one occurs then something has gone wrong with my logic
	AllBucketsOccupied HashTableError = HashTableError("all buckets occupied")

	ObjectDoesNotExistToUpdate HashTableError = HashTableError("object does not exist to update")

	ObjectDoesNotExistToDelete HashTableError = HashTableError("object does not exist to delete")

	BucketOccupied HashTableError = HashTableError("bucket is already occupied")
)
