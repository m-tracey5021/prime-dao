package ht

import (
	"encoding/binary"
	"os"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type MockHashable struct {
	MockHashableId uuid.UUID

	Data uint64
}

func (obj MockHashable) Id() uuid.UUID {

	return obj.MockHashableId
}

func (obj *MockHashable) SetId(id uuid.UUID) {

	obj.MockHashableId = id
}

func (obj *MockHashable) Compare(other schema.Orderable) schema.Order {

	return schema.Equal
}

func (obj MockHashable) Descriptor() string {

	return "mock_hash"
}

func (obj MockHashable) Size() int {

	return 24
}

func (obj MockHashable) WriteSelf(file *os.File) error {

	if err := binary.Write(file, binary.LittleEndian, obj.MockHashableId); err != nil {

		return err
	}
	if err := binary.Write(file, binary.LittleEndian, obj.Data); err != nil {

		return err
	}
	return nil
}

func (obj MockHashable) ReadSelf(file *os.File) (schema.FixedSize, error) {

	var mockHashableId uuid.UUID

	var data uint64

	if err := binary.Read(file, binary.BigEndian, &mockHashableId); err != nil {

		return nil, err
	}
	if err := binary.Read(file, binary.BigEndian, &data); err != nil {

		return nil, err
	}
	return MockHashable{mockHashableId, data}, nil
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
