package lib

import "slices"

func BoolToByte(value bool) byte {

	if value {

		return 1
	}
	return 0
}

func ByteToBool(value byte) bool {

	return value != 0
}

func NewId(slice []uint64) uint64 {

	var id uint64 = 0

	idExists := slices.Contains(slice, id)

	for idExists {

		id += 1

		idExists = slices.Contains(slice, id)
	}
	return id
}

func RemoveId(id uint64, slice []uint64) []uint64 {

	matches := func(element uint64) bool {

		return element == id
	}
	slice = slices.DeleteFunc(slice, matches)

	return slice
}
