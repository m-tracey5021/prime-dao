package dao

import (
	"bytes"
	"sync"

	"github.com/google/uuid"
	"github.com/m-tracey5021/prime-dao/pkg/bt"
	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/ht"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

// U is the sort key type

type Dao[T schema.Orderable] struct {
	fileManager fm.IFileManager

	metadata DaoMetadata[T]

	metadataIO dataio.DataIO[DaoMetadata[T]]

	indexHashTable ht.ITSFHashTable[DaoIndex]

	indexBTree bt.BTree[DaoIndex, T]

	objectIO dataio.DataIO[T]

	idMutex sync.Mutex

	cacheMutex sync.Mutex

	metadataMutex sync.Mutex
}

func From[T schema.Orderable](fileManager fm.IFileManager) (*Dao[T], error) {

	var described T

	fileManager = fileManager.Concatenate(described.Descriptor())

	daoMetadataFile, err := fileManager.OpenAndLock(fm.DaoMetadataFile)

	defer fileManager.CloseAndUnlock(daoMetadataFile, &err)

	if err != nil {

		return &Dao[T]{}, err
	}
	size, err := fileManager.Size(daoMetadataFile)

	if err != nil {

		return &Dao[T]{}, err
	}
	metadataIO := dataio.DataIO[DaoMetadata[T]]{}

	var metadata DaoMetadata[T]

	if size > 0 {

		metadata, err = metadataIO.ReadSizePrefixed(daoMetadataFile)

		if err != nil {

			return &Dao[T]{}, err
		}

	} else {

		initialTable := uint64(0)

		tableMapping := make(map[uint64]int)

		tableMapping[initialTable] = 0 // TODO do these fields need to be init'd or can they be iffed elsewhere in the logic

		metadata = DaoMetadata[T]{

			ObjectIdStore: NewIdStore(),

			TableIdStore: DaoIdStore{int(initialTable), make([]uint64, 0)},

			Cache: NewCache[T](),

			MaxObjects: uint64(2),

			AvailableTableObjectCount: tableMapping,
		}
		if _, err := metadataIO.WriteSizePrefixed(daoMetadataFile, metadata); err != nil {

			return &Dao[T]{}, err
		}
	}
	dao := Dao[T]{

		fileManager: fileManager,

		metadata: metadata,

		metadataIO: metadataIO,

		indexHashTable: ht.FromFileManager[DaoIndex](fileManager.Concatenate("idx")),

		indexBTree: bt.BTree[DaoIndex, T]{},

		objectIO: dataio.DataIO[T]{},
	}
	btree, err := bt.NewBTree[DaoIndex, T](3, fileManager, dao.GetCompareValue)
	/*
		TODO do this a bit cleaner, it is odd creating the dao, then
		setting the function on the btree from the dao, to then add
		it back to the dao again afterwards. This also means we have to
		return a pointer to the dao because it copies the lock
	*/

	if err != nil {

		return &Dao[T]{}, err
	}
	dao.indexBTree = btree

	return &dao, nil
}

func (dao *Dao[T]) NewTransaction() DaoTransaction[T] {

	return DaoTransaction[T]{

		requestIds: map[uint64]struct{}{},

		requests: []IProcessableRequest[T]{},

		dependencyResolver: NewResolver[T](),

		dependencies: []uint64{},
	}
}

func (dao *Dao[T]) ExecuteTransaction(transaction DaoTransaction[T]) (map[uint64]IResult[T], error) {

	mappedResults := make(map[uint64]IResult[T])

	transactionQueue := NewQueue[T](10, 10, 100)

	transactionQueue.Start(transaction.dependencyResolver, dao)

	transactionQueue.ProcessAsync(transaction.requests...)

	results := transactionQueue.Stop()

	daoMetadataFile, err := dao.fileManager.OpenAndLock(fm.DaoMetadataFile)

	defer dao.fileManager.CloseAndUnlock(daoMetadataFile, &err)

	_, err = dao.metadataIO.WriteSizePrefixed(daoMetadataFile, dao.metadata)

	if err != nil {

		return map[uint64]IResult[T]{}, err
	}
	for _, result := range results {

		mappedResults[result.RequestId()] = result.Result()
	}
	return mappedResults, nil
}

/*
	Can extend this functions at some point to be
	able to add different comparators for different types
	and build the btree based on that order
*/

func (dao *Dao[T]) Compare(a, b DaoIndex) int {

	objectA, err := dao.Get(a.id)

	if err != nil {

		return 0
	}
	objectB, err := dao.Get(b.id)

	uuidCompare := CompareUUIDs((*objectA).Id(), (*objectB).Id())

	return uuidCompare
}

func (dao *Dao[T]) GetCompareValue(indexId uuid.UUID) (*T, error) {

	return dao.Get(indexId)
}

func CompareUUIDs(a, b uuid.UUID) int {

	return bytes.Compare(a[:], b[:]) // Lexicographic comparison
}
