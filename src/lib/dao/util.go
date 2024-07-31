package dao

import "slices"

func boolToByte(value bool) byte {

	if value {

		return 1
	}
	return 0
}

func byteToBool(value byte) bool {

	return value != 0
}

func smallestMissing(slice []uint64) uint64 {

	var id uint64 = 0

	idExists := slices.Contains(slice, id)

	for idExists {

		id += 1

		idExists = slices.Contains(slice, id)
	}
	return id
}
