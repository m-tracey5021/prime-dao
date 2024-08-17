package ht

import (
	"encoding/binary"
	"os"
	"transformer/src/lib/dao/fm"
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"

	"github.com/stretchr/testify/mock"
)

// import (
// 	"encoding/binary"
// 	"os"
// 	"transformer/src/lib/dao/fm"
// 	"transformer/src/lib/dao/schema"
// 	"transformer/src/lib/daoio"
// 	"unsafe"

// 	"github.com/stretchr/testify/mock"
// )

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

func setupMockDao() (

	*fm.MockFileManager,

	*MockHashTableManager[MockHashable],

	*daoio.MockFixedSizeDao[HashTableBucketHeader],

	*daoio.MockFixedSizeDao[MockHashable],

	TSFHashTable[MockHashable],

) {
	id := uint64(0)

	mockFileManager := new(fm.MockFileManager)

	mockTableManager := new(MockHashTableManager[MockHashable])

	mockBucketHeaderIO := new(daoio.MockFixedSizeDao[HashTableBucketHeader])

	mockObjectIO := new(daoio.MockFixedSizeDao[MockHashable])

	hashTable := Inject[MockHashable](id, mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO)

	return mockFileManager, mockTableManager, mockBucketHeaderIO, mockObjectIO, hashTable
}

func setupMockTableManagerDependencies() (

	*fm.MockFileManager,

	*daoio.MockFixedSizeDao[HashTableBucketHeader],

	*daoio.MockFixedSizeDao[MockHashable],

	HashTableManager[MockHashable],

) {

	mockFileManager := new(fm.MockFileManager)

	mockBucketHeaderIO := new(daoio.MockFixedSizeDao[HashTableBucketHeader])

	mockObjectIO := new(daoio.MockFixedSizeDao[MockHashable])

	tableManager := InjectTableManager(mockFileManager, mockBucketHeaderIO, mockObjectIO)

	return mockFileManager, mockBucketHeaderIO, mockObjectIO, tableManager
}

func setupMockFiles(mockFileManager *fm.MockFileManager) (*os.File, *os.File, *os.File) {

	mockManagingFile := new(os.File)

	mockMainTable := new(os.File)

	mockCollisionTable := new(os.File)

	mockFileManager.On("Close", mock.AnythingOfType("*os.File")).Return(nil)

	mockFileManager.On("OpenAndLock", fm.HashTableManagingFile, mock.AnythingOfType("uint64")).Return(mockManagingFile, nil)

	mockFileManager.On("OpenAndLock", fm.HashTableMainTable, mock.AnythingOfType("uint64")).Return(mockMainTable, nil)

	mockFileManager.On("OpenAndLock", fm.HashTableCollisionTable, mock.AnythingOfType("uint64")).Return(mockCollisionTable, nil)

	return mockManagingFile, mockMainTable, mockCollisionTable
}

// // func setupMockDao() (

// // 	*MockDaoFileContainer,

// // 	*MockDaoIO[DaoIdentifierCache],

// // 	*MockFixedSizeDao[DaoBucketHeader],

// // 	*MockFixedSizeDao[MockHashable],

// // 	*FixedSizeDao[MockHashable],

// // 	error,

// // ) {

// // 	mockFileManager := new(MockDaoFileContainer)

// // 	mockCacheIO := new(MockDaoIO[DaoIdentifierCache])

// // 	mockBucketHeaderIO := new(MockFixedSizeDao[DaoBucketHeader])

// // 	mockObjectIO := new(MockFixedSizeDao[MockHashable])

// // 	mockManagingFile := new(os.File)

// // mockFileManager.On("ManagingFile").Return(mockManagingFile, nil)

// // mockFileManager.On("Size", mockManagingFile).Return(0, nil)

// // mockFileManager.On("Close", mockManagingFile).Return(nil)

// // 	mockCacheIO.On("WriteSizePrefixed", mockManagingFile, mock.Anything).Return(0, nil)

// // 	dao, err := NewDao(mockFileManager, mockCacheIO, mockBucketHeaderIO, mockObjectIO)

// // 	if err != nil {

// // 		return nil, nil, nil, nil, nil, err
// // 	}
// // 	return mockFileManager, mockCacheIO, mockBucketHeaderIO, mockObjectIO, dao, nil
// // }
