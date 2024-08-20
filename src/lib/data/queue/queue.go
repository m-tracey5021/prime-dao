package queue

import (
	"sync"
)

type TSFQueueResult[T any] struct {
	request IRequest[T]

	result IResult[T]
}

func (result TSFQueueResult[T]) Result() IResult[T] {

	return result.result
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

func (queue *TSFQueue[T]) NewRequestId() uint64 {

	return 0

	// id := lib.NewId(queue.requestIds)

	// queue.mu.Lock()

	// queue.requestIds = append(queue.requestIds, id)

	// queue.mu.Unlock()

	// return id
}

func (queue *TSFQueue[T]) Start() {

	for i := 0; i < queue.numberOfWorkers; i++ {

		queue.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer queue.requestsWaitGroup.Done()

			for batch := range queue.channel {

				for _, request := range batch {

					result := queue.requestProcessor.Process(request)

					queue.results <- TSFQueueResult[T]{request, result}
				}
			}
		}(i)
	}
	queue.resultsWaitGroup.Add(1)

	go func() {

		defer queue.resultsWaitGroup.Done()

		for result := range queue.results {

			queue.mu.Lock()

			queue.requestsProcessed = append(queue.requestsProcessed, result)

			queue.mu.Unlock()
		}
	}()
}

func (queue *TSFQueue[T]) ProcessSync(request IRequest[T]) IResult[T] {

	return queue.requestProcessor.Process(request)
}

func (queue *TSFQueue[T]) ProcessAsync(requests ...IRequest[T]) {

	numberOfBatches := (len(requests) + queue.batchSize - 1) / queue.batchSize

	for i := 0; i < numberOfBatches; i++ {

		start := i * queue.batchSize

		end := (i + 1) * queue.batchSize

		if end > len(requests) {

			end = len(requests)
		}
		batch := requests[start:end]

		queue.channel <- batch
	}
}

func (queue *TSFQueue[T]) Stop() []TSFQueueResult[T] {

	close(queue.channel)

	queue.requestsWaitGroup.Wait()

	close(queue.results)

	queue.resultsWaitGroup.Wait()

	return queue.requestsProcessed
}
