package queue

import "sync"

type TSFProcessResult[T any] struct {
	request IProcessableRequest[T]

	result IResult[T]
}

type TSFQueue[T any] struct {
	numberOfWorkers int

	batchSize int

	bufferSize int

	channel chan []IProcessableRequest[T]

	results chan TSFProcessResult[T]

	requestIds []uint64

	requestsProcessed []TSFProcessResult[T]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex
}

func NewQueueB[T any](numberOfWorkers int, batchSize int, bufferSize int) *TSFQueue[T] {

	return &TSFQueue[T]{

		numberOfWorkers: numberOfWorkers,

		batchSize: batchSize,

		bufferSize: bufferSize,

		channel: make(chan []IProcessableRequest[T], bufferSize),

		results: make(chan TSFProcessResult[T], bufferSize),

		requestIds: make([]uint64, 0),

		requestsProcessed: make([]TSFProcessResult[T], 0),
	}
}

func (queue *TSFQueue[T]) Start() {

	for i := 0; i < queue.numberOfWorkers; i++ {

		queue.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer queue.requestsWaitGroup.Done()

			for batch := range queue.channel {

				for _, request := range batch {

					result := request.Process()

					queue.results <- TSFProcessResult[T]{request, result}
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

func (queue *TSFQueue[T]) ProcessSync(request IProcessableRequest[T]) IResult[T] {

	return request.Process()
}

func (queue *TSFQueue[T]) ProcessAsync(requests ...IProcessableRequest[T]) {

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

func (queue *TSFQueue[T]) Stop() []TSFProcessResult[T] {

	close(queue.channel)

	queue.requestsWaitGroup.Wait()

	close(queue.results)

	queue.resultsWaitGroup.Wait()

	return queue.requestsProcessed
}
