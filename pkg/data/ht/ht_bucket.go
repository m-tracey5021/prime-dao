package ht

import "prime-dao/pkg/data/schema"

type HashTableBucket[T schema.FixedSizeIdentifiable] struct {
	bucketLocation int

	objectLocation int

	bucketHeader HashTableBucketHeader

	object T
}
