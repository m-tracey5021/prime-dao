package dao

import (
	"sync"

	"github.com/m-tracey5021/prime-dao/pkg/dataio"
	"github.com/m-tracey5021/prime-dao/pkg/fm"
	"github.com/m-tracey5021/prime-dao/pkg/ht"
	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type Dao[T schema.Identifiable] struct {
	fileManager fm.IFileManager

	metadata DaoMetadata[T]

	metadataIO dataio.DataIO[DaoMetadata[T]]

	indexHashTable ht.ITSFHashTable[DaoIndex]

	objectIO dataio.DataIO[T]

	idMutex sync.Mutex

	cacheMutex sync.Mutex

	metadataMutex sync.Mutex
}

func From[T schema.Identifiable](fileManager fm.IFileManager) (Dao[T], error) {

	var described T

	fileManager = fileManager.Concatenate(described.Descriptor())

	daoMetadataFile, err := fileManager.OpenAndLock(fm.DaoMetadataFile)

	defer fileManager.CloseAndUnlock(daoMetadataFile, &err)

	if err != nil {

		return Dao[T]{}, err
	}
	size, err := fileManager.Size(daoMetadataFile)

	if err != nil {

		return Dao[T]{}, err
	}
	metadataIO := dataio.DataIO[DaoMetadata[T]]{}

	var metadata DaoMetadata[T]

	if size > 0 {

		metadata, err = metadataIO.ReadSizePrefixed(daoMetadataFile)

		if err != nil {

			return Dao[T]{}, err
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

			return Dao[T]{}, err
		}
	}
	return Dao[T]{

			fileManager: fileManager,

			metadata: metadata,

			metadataIO: metadataIO,

			indexHashTable: ht.FromFileManager[DaoIndex](fileManager.Concatenate("idx")),

			objectIO: dataio.DataIO[T]{},
		},
		err
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
