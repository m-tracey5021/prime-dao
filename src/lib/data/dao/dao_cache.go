package dao

import (
	"encoding/binary"
	"os"
	"sync"
	"transformer/src/lib/data"
	"transformer/src/lib/data/schema"
)

type DaoCache struct {
	lastAddedExists bool

	lastAddedId uint64

	lastDeletedExists bool

	lastDeletedId uint64
}

func (cache *DaoCache) NewId(mu *sync.Mutex) uint64 {

	mu.Lock()

	if cache.lastDeletedExists {

		return cache.lastDeletedId
	}
	if cache.lastAddedExists {

		cache.lastAddedId += 1

		return cache.lastAddedId
	}
	cache.lastAddedExists = true

	mu.Unlock()

	return uint64(0)
}

func (cache *DaoCache) DeleteId(id uint64, mu *sync.Mutex) {

	mu.Lock()

	cache.lastDeletedExists = true

	cache.lastDeletedId = id

	mu.Unlock()
}

func (cache DaoCache) Size() int {

	return 18
}

func (cache DaoCache) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, cache.lastAddedExists); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, cache.lastAddedId); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, cache.lastDeletedExists); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, cache.lastDeletedId); err != nil {

		return err
	}
	return nil
}

func (cache DaoCache) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var lastAddedExists byte

	var lastAddedId uint64

	var lastDeletedExists byte

	var lastDeletedId uint64

	if err := binary.Read(file, binary.LittleEndian, &lastAddedExists); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &lastAddedId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &lastDeletedExists); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.LittleEndian, &lastDeletedId); err != nil {

		return nil, err
	}
	return DaoCache{data.ByteToBool(lastAddedExists), lastAddedId, data.ByteToBool(lastDeletedExists), lastDeletedId}, nil
}
