package ht

import (
	"errors"
	"io"
	"transformer/src/lib/dao/fm"
)

func (ht *TSFHashTable[T]) SaveNewCollisionTableId(id uint64) error {

	managingFile, err := ht.fileContainer.Open(fm.HashTableManagingFile, ht.id)

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := ht.fileContainer.Close(managingFile); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	ht.identifierCache.CollisionTableIds = append(ht.identifierCache.CollisionTableIds, id)

	_, err = ht.cacheIO.WriteSizePrefixed(managingFile, ht.identifierCache)

	if err != nil {

		return err
	}
	return err
}

func (ht *TSFHashTable[T]) Save(object T) error {

	mainTable, err := ht.fileContainer.Open(fm.HashTableMainTable, ht.id)

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := ht.fileContainer.Close(mainTable); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	position := ht.hash(object.Id())

	if err := ht.fileContainer.GoTo(position, mainTable); err != nil {

		return err
	}
	bucketHeader, err := ht.bucketHeaderIO.Read(mainTable)

	if err != nil {

		if errors.Is(err, io.EOF) {

			err = nil

			bucketHeader = HashTableBucketHeader{false, false, 0}

		} else {

			return err
		}
	}
	if bucketHeader.occupied {

		readObject, err := ht.objectIO.Read(mainTable)

		if err != nil {

			return err
		}
		if object.Id() == readObject.Id() {

			return ObjectAlreadyExists
		}
		if !bucketHeader.previousCollision {

			collisionTableId := ht.NewCollisionTableId()

			if err := ht.SaveForCollision(object, collisionTableId); err != nil {

				return err
			}
			if err := ht.SaveNewCollisionTableId(collisionTableId); err != nil {

				return err
			}
			updatedBucketHeader := HashTableBucketHeader{true, true, collisionTableId}

			if err := ht.fileContainer.GoTo(position, mainTable); err != nil {

				return err
			}
			if err := ht.bucketHeaderIO.Write(mainTable, updatedBucketHeader); err != nil {

				return nil
			}

		} else {

			if err := ht.SaveForCollision(object, bucketHeader.collisionTableId); err != nil {

				return err
			}
		}

	} else {

		bucketHeader = HashTableBucketHeader{true, false, 0}

		if err := ht.fileContainer.GoTo(position, mainTable); err != nil {

			return err
		}
		if err := ht.bucketHeaderIO.Write(mainTable, bucketHeader); err != nil {

			return err
		}
		if err := ht.objectIO.Write(mainTable, object); err != nil {

			return err
		}
	}
	return err
}

// func (ht *TSFHashTable[T]) SaveForCollision(object T, initialHash int, collisionTableId uint64) error {

// 	collisionNumber, collisionTableIdNew := ht.hashCollision(object.Id())

// 	collisionTableSuffix := fmt.Sprintf("%v_%v", initialHash, collisionTableIdNew)

// 	collisionTable, err := ht.fileContainer.Open(fm.HashTableCollisionTable, collisionTableSuffix)

// 	if err != nil {

// 		return err
// 	}
// 	defer func() {

// 		if innerErr := ht.fileContainer.Close(collisionTable); innerErr != nil {

// 			if err != nil {

// 				err = errors.Join(innerErr, err)

// 			} else {

// 				err = innerErr
// 			}
// 		}
// 	}()

// 	for {

// 		readObject, err := ht.objectIO.Read(collisionTable)

// 		if err != nil {

// 			if errors.Is(err, io.EOF) {

// 				if err := ht.objectIO.Write(collisionTable, object); err != nil {

// 					return err
// 				}
// 				err = nil

// 				return err
// 			}
// 			return err
// 		}
// 		if object.Id() == readObject.Id() {

// 			return ObjectAlreadyExists
// 		}
// 	}
// }

func (ht *TSFHashTable[T]) SaveForCollision(object T, initialHash int) error {

	collisionPosition, collisionTableId := ht.hashCollision(object.Id(), initialHash)

	collisionTable, err := ht.fileContainer.Open(fm.HashTableCollisionTable, collisionTableId)

	if err != nil {

		return err
	}
	defer func() {

		if innerErr := ht.fileContainer.Close(collisionTable); innerErr != nil {

			if err != nil {

				err = errors.Join(innerErr, err)

			} else {

				err = innerErr
			}
		}
	}()

	if err := ht.fileContainer.GoTo(collisionPosition, collisionTable); err != nil {

		return err
	}
	bucketHeader, err := ht.bucketHeaderIO.Read(collisionTable)

	if bucketHeader.occupied {

		readObject, err := ht.objectIO.Read(collisionTable)

		if err != nil {

			return err
		}
		if object.Id() == readObject.Id() {

			return ObjectAlreadyExists
		}
		return CollisionAlreadyOccupied
	}
	if err := ht.objectIO.Write(collisionTable, object); err != nil {

		return err
	}
	return err
}
