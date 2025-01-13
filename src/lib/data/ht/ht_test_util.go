package ht

import (
	"encoding/binary"
	"os"
	"prime-dao/src/lib/data/dataio"
	"prime-dao/src/lib/data/fm"
	"prime-dao/src/lib/data/schema"
)

type MockHashable struct {
	MockHashableId uint64
}

func (obj MockHashable) Id() uint64 {

	return obj.MockHashableId
}

func (obj MockHashable) SetId(id uint64) schema.Identifiable {

	return MockHashable{id}
}

func (obj MockHashable) Size() int {

	return 8
}

func (obj MockHashable) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, obj.MockHashableId); err != nil {

		return err
	}
	return nil
}

func (obj MockHashable) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var mockHashableId uint64

	if err := binary.Read(file, binary.LittleEndian, &mockHashableId); err != nil {

		return nil, err
	}
	return MockHashable{mockHashableId}, nil
}

func setupHashTableMockDependencies() (

	*fm.MockFileManager,

	*MockHashTableManager[MockHashable],

	*dataio.MockFixedSizeDataIO[HashTableBucketHeader],

	*dataio.MockFixedSizeDataIO[MockHashable],

	TSFHashTable[MockHashable],

) {
	id := uint64(0)

	mockFileManager := new(fm.MockFileManager)

	mockTableManager := new(MockHashTableManager[MockHashable])

	mockBucketHeaderIO := new(dataio.MockFixedSizeDataIO[HashTableBucketHeader])

	mockObjectIO := new(dataio.MockFixedSizeDataIO[MockHashable])

	hashTable := TSFHashTable[MockHashable]{id, mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO}

	return mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, hashTable
}

func setupTableManagerMockDependencies() (

	*fm.MockFileManager,

	*dataio.MockFixedSizeDataIO[HashTableBucketHeader],

	*dataio.MockFixedSizeDataIO[MockHashable],

	HashTableManager[MockHashable],

) {

	mockFileManager := new(fm.MockFileManager)

	mockBucketHeaderIO := new(dataio.MockFixedSizeDataIO[HashTableBucketHeader])

	mockObjectIO := new(dataio.MockFixedSizeDataIO[MockHashable])

	tableManager := HashTableManager[MockHashable]{10, 10, 10, mockFileManager, mockBucketHeaderIO, mockObjectIO}

	return mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager
}
