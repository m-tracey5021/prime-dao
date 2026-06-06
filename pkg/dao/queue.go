package dao

import (
	"fmt"
	"sync"

	"github.com/m-tracey5021/prime-dao/pkg/schema"
	"golang.org/x/exp/constraints"
)

type TSFQueueResult[T schema.Orderable[U], U constraints.Ordered] struct {
	requestId uint64

	result IResult[T]
}

func (queueResult TSFQueueResult[T, U]) RequestId() uint64 {

	return queueResult.requestId
}

func (queueResult TSFQueueResult[T, U]) Result() IResult[T] {

	return queueResult.result
}

type Queue[T schema.Orderable[U], U constraints.Ordered] struct {
	numberOfWorkers int

	batchSize int

	bufferSize int

	channel chan []IProcessableRequest[T, U]

	results chan TSFQueueResult[T, U]

	requestsProcessed []TSFQueueResult[T, U]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex
}

func NewQueue[T schema.Orderable[U], U constraints.Ordered](numberOfWorkers int, batchSize int, bufferSize int) *Queue[T, U] {

	return &Queue[T, U]{

		numberOfWorkers: numberOfWorkers,

		batchSize: batchSize,

		bufferSize: bufferSize,

		channel: make(chan []IProcessableRequest[T, U], bufferSize),

		results: make(chan TSFQueueResult[T, U], bufferSize),

		requestsProcessed: make([]TSFQueueResult[T, U], 0),
	}
}

func (queue *Queue[T, U]) Start(dependencyResolver *QueueDependencyResolver[T, U], dao *Dao[T, U]) {

	dependencyResolver.ListenForCompletedDependencies()

	for i := 0; i < queue.numberOfWorkers; i++ {

		queue.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer queue.requestsWaitGroup.Done()

			for batch := range queue.channel {

				for _, request := range batch {

					queue.ProcessDependencies(request)

					result := request.Process(dao)

					mappedResult := TSFQueueResult[T, U]{request.RequestId(), result}

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

func (queue *Queue[T, U]) ProcessDependencies(processor IProcessableRequest[T, U]) {

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

func (queue *Queue[T, U]) ProcessAsync(requests ...IProcessableRequest[T, U]) {

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

func (queue *Queue[T, U]) Stop() []TSFQueueResult[T, U] {

	close(queue.channel)

	queue.requestsWaitGroup.Wait()

	close(queue.results)

	queue.resultsWaitGroup.Wait()

	return queue.requestsProcessed
}
