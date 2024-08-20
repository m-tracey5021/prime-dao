package req

import (
	"transformer/src/lib/data/dao"
	"transformer/src/lib/data/schema"
)

type DaoRequestProcessor[T schema.Identifiable] struct {
	dao dao.TSFDao[T]
}

// func (processor DaoRequestProcessor[T]) Process(request queue.IRequest[T]) queue.IResult[T] {

// 	switch request.Type() {

// 	case queue.SaveRequest:

// 		return processor.ProcessSaveRequest(request)

// 	case queue.GetRequest:

// 		return processor.ProcessGetRequest(request)

// 	case queue.UpdateRequest:

// 		return processor.ProcessUpdateRequest(request)

// 	case queue.DeleteRequest:

// 		return processor.ProcessDeleteRequest(request)

// 	default:

// 		return nil
// 	}
// }
