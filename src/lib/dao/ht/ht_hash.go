package ht

type Hash struct {
	tableGroup int

	tableNumber int

	position int
}

type HashMetrics struct {
	bucketSize int

	tableSize int

	maxCollisions int
}

func (metrics HashMetrics) ComputeHash(id uint64) Hash {

	// Calculate the table group to store the data in
	hash := int(id) % metrics.tableSize

	// Calculate the order of the hash i.e. where it sits in relation to the others if hashed
	order := int(id) / metrics.tableSize

	// Calculate the file/partition in which the data is stored
	tableNumber := order / metrics.maxCollisions

	// Calculate the actual position in the file based on bucket size and table number
	position := (order - (metrics.maxCollisions * tableNumber)) * metrics.bucketSize

	return Hash{hash, tableNumber, position}
}
