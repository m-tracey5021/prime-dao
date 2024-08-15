package ht

import (
	"transformer/src/lib/dao/schema"
	"transformer/src/lib/daoio"
)

func (processor HashTableProcessor[T]) ProcessDeleteRequest(request IRequest[T]) IResult[T] {

	table, bucket, err := processor.hashManager.Locate(request.ObjectId())

	defer processor.fileManager.CloseAndUnlock(table, &err)

	if err != nil {

		return &HashTableDeleteResult[T]{err: err}
	}
	if err := processor.fileManager.GoTo(bucket.bucketLocation, table); err != nil {

		return &HashTableDeleteResult[T]{err: err}
	}
	bucketHeader := HashTableBucketHeader{false, true, 0}

	if err := daoio.Write(table, bucketHeader); err != nil {

		return &HashTableDeleteResult[T]{err: err}
	}
	if err := daoio.Zero[T](table); err != nil {

		return &HashTableDeleteResult[T]{err: err}
	}
	return &HashTableDeleteResult[T]{err: err}
}

type HashTableDeleteRequest[T schema.FixedSizeIdentifiable] struct {
	requestId uint64

	objectId uint64
}

func (request *HashTableDeleteRequest[T]) Type() RequestType {

	return DeleteRequest
}

func (request *HashTableDeleteRequest[T]) RequestId() uint64 {

	return request.requestId
}

func (request *HashTableDeleteRequest[T]) ObjectId() uint64 {

	return request.objectId
}

func (request *HashTableDeleteRequest[T]) Object() *T {

	return nil
}

type HashTableDeleteResult[T schema.FixedSizeIdentifiable] struct {
	err error
}

func (result *HashTableDeleteResult[T]) Object() *T {

	return nil
}

func (result *HashTableDeleteResult[T]) Error() error {

	return result.err
}
