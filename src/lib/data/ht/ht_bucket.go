package ht

import "tsf-dao/src/lib/data/schema"

type HashTableBucket[T schema.FixedSizeIdentifiable] struct {
	bucketLocation int

	objectLocation int

	bucketHeader HashTableBucketHeader

	object T
}
