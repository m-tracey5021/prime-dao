package ht

import (
	"encoding/binary"
	"os"

	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
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
	mockFileManager := new(fm.MockFileManager)

	mockTableManager := new(MockHashTableManager[MockHashable])

	mockBucketHeaderIO := new(dataio.MockFixedSizeDataIO[HashTableBucketHeader])

	mockObjectIO := new(dataio.MockFixedSizeDataIO[MockHashable])

	hashTable := TSFHashTable[MockHashable]{mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO}

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
