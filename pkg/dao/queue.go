package dao

import (
	"fmt"
	"sync"

	"github.com/m-tracey5021/prime-dao/pkg/schema"
)

type TSFQueueResult[T schema.Orderable] struct {
	requestId uint64

	result IResult[T]
}

func (queueResult TSFQueueResult[T]) RequestId() uint64 {

	return queueResult.requestId
}

func (queueResult TSFQueueResult[T]) Result() IResult[T] {

	return queueResult.result
}

type Queue[T schema.Orderable] struct {
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

func NewQueue[T schema.Orderable](numberOfWorkers int, batchSize int, bufferSize int) *Queue[T] {

	return &Queue[T]{

		numberOfWorkers: numberOfWorkers,

		batchSize: batchSize,

		bufferSize: bufferSize,

		channel: make(chan []IProcessableRequest[T], bufferSize),

		results: make(chan TSFQueueResult[T], bufferSize),

		requestsProcessed: make([]TSFQueueResult[T], 0),
	}
}

func (queue *Queue[T]) Start(dependencyResolver *QueueDependencyResolver[T], dao *Dao[T]) {

	dependencyResolver.ListenForCompletedDependencies()

	for i := 0; i < queue.numberOfWorkers; i++ {

		queue.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer queue.requestsWaitGroup.Done()

			for batch := range queue.channel {

				for _, request := range batch {

					queue.ProcessDependencies(request)

					result := request.Process(dao)

					mappedResult := TSFQueueResult[T]{request.RequestId(), result}

					dependencyResolver.completed <- request.RequestId()

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

func (queue *Queue[T]) ProcessDependencies(processor IProcessableRequest[T]) {

	var dependencyWaitGroup sync.WaitGroup

	dependencyWaitGroup.Add(processor.GetNumDeps())

	go func() {

		for completed := range processor.Dependencies() {

			dependencyWaitGroup.Done()

			fmt.Printf("\nrequest %v waited for dependency %v", processor.RequestId(), completed)
		}
	}()
	dependencyWaitGroup.Wait()

	close(processor.Dependencies())
}

func (queue *Queue[T]) ProcessAsync(requests ...IProcessableRequest[T]) {

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

func (queue *Queue[T]) Stop() []TSFQueueResult[T] {

	close(queue.channel)

	queue.requestsWaitGroup.Wait()

	close(queue.results)

	queue.resultsWaitGroup.Wait()

	return queue.requestsProcessed
}
