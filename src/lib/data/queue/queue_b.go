package queue

import "sync"

type TSFProcessResult[T any] struct {
	request IProcessor[T]

	result IResult[T]
}

type TSFQueueB[T any] struct {
	numberOfWorkers int

	batchSize int

	bufferSize int

	channel chan []IProcessor[T]

	results chan TSFProcessResult[T]

	requestIds []uint64

	requestsProcessed []TSFProcessResult[T]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex
}

func NewQueueB[T any](numberOfWorkers int, batchSize int, bufferSize int) *TSFQueueB[T] {

	return &TSFQueueB[T]{

		numberOfWorkers: numberOfWorkers,

		batchSize: batchSize,

		bufferSize: bufferSize,

		channel: make(chan []IProcessor[T], bufferSize),

		results: make(chan TSFProcessResult[T], bufferSize),

		requestIds: make([]uint64, 0),

		requestsProcessed: make([]TSFProcessResult[T], 0),
	}
}

func (queue *TSFQueueB[T]) NewRequestId() uint64 {

	return 0

	// id := lib.NewId(queue.requestIds)

	// queue.mu.Lock()

	// queue.requestIds = append(queue.requestIds, id)

	// queue.mu.Unlock()

	// return id
}

func (queue *TSFQueueB[T]) Start() {

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

func (queue *TSFQueueB[T]) ProcessSync(request IProcessor[T]) IResult[T] {

	return request.Process()
}

func (queue *TSFQueueB[T]) ProcessAsync(requests ...IProcessor[T]) {

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

func (queue *TSFQueueB[T]) Stop() []TSFProcessResult[T] {

	close(queue.channel)

	queue.requestsWaitGroup.Wait()

	close(queue.results)

	queue.resultsWaitGroup.Wait()

	return queue.requestsProcessed
}
