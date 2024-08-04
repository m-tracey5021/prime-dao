package dao

import (
	"encoding/binary"
	"os"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
	"unsafe"

	"github.com/stretchr/testify/mock"
)

type MockHashable struct {
	MockHashableId uint64
}

func (obj MockHashable) Id() uint64 {

	return obj.MockHashableId
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

func setupMockDao() (

	*MockFileManager,

	*daoio.MockDaoIO[HashTableIdentifierCache],

	*daoio.MockFixedSizeDao[HashTableBucketHeader],

	*daoio.MockFixedSizeDao[MockHashable],

	TSFHashTable[MockHashable],

) {
	id := uint64(0)

	mockFileContainer := new(MockFileManager)

	mockCacheIO := new(daoio.MockDaoIO[HashTableIdentifierCache])

	mockBucketHeaderIO := new(daoio.MockFixedSizeDao[HashTableBucketHeader])

	mockObjectIO := new(daoio.MockFixedSizeDao[MockHashable])

	bucketSize := int(unsafe.Sizeof(*new(HashTableBucketHeader)) + unsafe.Sizeof(*new(MockHashable)))

	tableSize := 10

	collisionTableIds := make([]uint64, 0)

	identifierCache := HashTableIdentifierCache{collisionTableIds}

	dao := Default[MockHashable](id, bucketSize, tableSize, mockFileContainer, identifierCache, mockCacheIO, mockBucketHeaderIO, mockObjectIO)

	return mockFileContainer, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao
}

func setupMockFiles(mockFileContainer *MockFileManager) (*os.File, *os.File, *os.File) {

	mockManagingFile := new(os.File)

	mockMainTable := new(os.File)

	mockCollisionTable := new(os.File)

	mockFileContainer.On("Close", mock.AnythingOfType("*os.File")).Return(nil)

	mockFileContainer.On("File", HashTableManagingFile, mock.AnythingOfType("uint64")).Return(mockManagingFile, nil)

	mockFileContainer.On("File", HashTableMainTable, mock.AnythingOfType("uint64")).Return(mockMainTable, nil)

	mockFileContainer.On("File", HashTableCollisionTable, mock.AnythingOfType("uint64")).Return(mockCollisionTable, nil)

	return mockManagingFile, mockMainTable, mockCollisionTable
}

// func setupMockDao() (

// 	*MockDaoFileContainer,

// 	*MockDaoIO[DaoIdentifierCache],

// 	*MockFixedSizeDao[DaoBucketHeader],

// 	*MockFixedSizeDao[MockHashable],

// 	*FixedSizeDao[MockHashable],

// 	error,

// ) {

// 	mockFileContainer := new(MockDaoFileContainer)

// 	mockCacheIO := new(MockDaoIO[DaoIdentifierCache])

// 	mockBucketHeaderIO := new(MockFixedSizeDao[DaoBucketHeader])

// 	mockObjectIO := new(MockFixedSizeDao[MockHashable])

// 	mockManagingFile := new(os.File)

// mockFileContainer.On("ManagingFile").Return(mockManagingFile, nil)

// mockFileContainer.On("Size", mockManagingFile).Return(0, nil)

// mockFileContainer.On("Close", mockManagingFile).Return(nil)

// 	mockCacheIO.On("WriteSizePrefixed", mockManagingFile, mock.Anything).Return(0, nil)

// 	dao, err := NewDao(mockFileContainer, mockCacheIO, mockBucketHeaderIO, mockObjectIO)

// 	if err != nil {

// 		return nil, nil, nil, nil, nil, err
// 	}
// 	return mockFileContainer, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao, nil
// }
