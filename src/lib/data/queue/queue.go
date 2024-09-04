package queue

import "sync"

type TSFQueueResult[T any] struct {
	requestId uint64

	result IResult[T]
}

func (queueResult TSFQueueResult[T]) RequestId() uint64 {

	return queueResult.requestId
}

func (queueResult TSFQueueResult[T]) Result() IResult[T] {

	return queueResult.result
}

type TSFQueue[T any] struct {
	numberOfWorkers int

	batchSize int

	bufferSize int

	channel chan []IProcessableRequest[T]

	results chan TSFQueueResult[T]

	requestsProcessed []TSFQueueResult[T]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex
}

func NewQueue[T any](numberOfWorkers int, batchSize int, bufferSize int) *TSFQueue[T] {

	return &TSFQueue[T]{

		numberOfWorkers: numberOfWorkers,

		batchSize: batchSize,

		bufferSize: bufferSize,

		channel: make(chan []IProcessableRequest[T], bufferSize),

		results: make(chan TSFQueueResult[T], bufferSize),

		requestsProcessed: make([]TSFQueueResult[T], 0),
	}
}

func (queue *TSFQueue[T]) Start(dependencyResolver *QueueDependencyResolver[T]) {

	dependencyResolver.ListenForCompletedDependencies()

	for i := 0; i < queue.numberOfWorkers; i++ {

		queue.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer queue.requestsWaitGroup.Done()

			for batch := range queue.channel {

				for _, request := range batch {

					var processWaitGroup sync.WaitGroup

					processWaitGroup.Add(1)

					result := request.ProcessWithDependencies(&processWaitGroup, dependencyResolver)

					mappedResult := TSFQueueResult[T]{request.RequestId(), result}

					requestContext := QueueDependency{request.RequestId(), &processWaitGroup}

					dependencyResolver.completed <- requestContext

					queue.results <- mappedResult
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

func (queue *TSFQueue[T]) Stop() []TSFQueueResult[T] {

	close(queue.channel)

	queue.requestsWaitGroup.Wait()

	close(queue.results)

	queue.resultsWaitGroup.Wait()

	return queue.requestsProcessed
}
