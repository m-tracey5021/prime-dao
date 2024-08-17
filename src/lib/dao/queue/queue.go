package queue

import (
	"sync"
)

type TSFQueueResult[T any] struct {
	request IRequest[T]

	result IResult[T]
}

type TSFQueue[T any] struct {
	numberOfWorkers int

	batchSize int

	bufferSize int

	requestProcessor IQueueProcessor[T]

	channel chan []IRequest[T]

	results chan TSFQueueResult[T]

	requestIds []uint64

	requestsProcessed []TSFQueueResult[T]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex
}

func NewQueue[T any](numberOfWorkers int, batchSize int, bufferSize int, requestProcessor IQueueProcessor[T]) *TSFQueue[T] {

	return &TSFQueue[T]{

		numberOfWorkers:   numberOfWorkers,
		batchSize:         batchSize,
		bufferSize:        bufferSize,
		requestProcessor:  requestProcessor,
		channel:           make(chan []IRequest[T], bufferSize),
		results:           make(chan TSFQueueResult[T], bufferSize),
		requestIds:        make([]uint64, 0),
		requestsProcessed: make([]TSFQueueResult[T], 0),
	}
}

func (cq *TSFQueue[T]) NewRequestId() uint64 {

	return 0

	// id := lib.NewId(cq.requestIds)

	// cq.mu.Lock()

	// cq.requestIds = append(cq.requestIds, id)

	// cq.mu.Unlock()

	// return id
}

func (cq *TSFQueue[T]) Start() {

	for i := 0; i < cq.numberOfWorkers; i++ {

		cq.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer cq.requestsWaitGroup.Done()

			for batch := range cq.channel {

				for _, request := range batch {

					result := cq.requestProcessor.Process(request)

					cq.results <- TSFQueueResult[T]{request, result}
				}
			}
		}(i)
	}
	cq.resultsWaitGroup.Add(1)

	go func() {

		defer cq.resultsWaitGroup.Done()

		for result := range cq.results {

			cq.mu.Lock()

			cq.requestsProcessed = append(cq.requestsProcessed, result)

			cq.mu.Unlock()
		}
	}()
}

func (cq *TSFQueue[T]) ProcessSync(request IRequest[T]) IResult[T] {

	return cq.requestProcessor.Process(request)
}

func (cq *TSFQueue[T]) ProcessAsync(requests ...IRequest[T]) {

	numberOfBatches := (len(requests) + cq.batchSize - 1) / cq.batchSize

	for i := 0; i < numberOfBatches; i++ {

		start := i * cq.batchSize

		end := (i + 1) * cq.batchSize

		if end > len(requests) {

			end = len(requests)
		}
		batch := requests[start:end]

		cq.channel <- batch
	}
}

func (cq *TSFQueue[T]) Stop() []TSFQueueResult[T] {

	close(cq.channel)

	cq.requestsWaitGroup.Wait()

	close(cq.results)

	cq.resultsWaitGroup.Wait()

	return cq.requestsProcessed
}
