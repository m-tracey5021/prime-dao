package ht

import (
	"sync"
	"transformer/src/lib/dao/schema"
)

type CommandQueueResult[T any] struct {
	request IRequest[T]

	result IResult[T]
}

type CommandQueue[T any] struct {
	numberOfWorkers int

	batchSize int

	bufferSize int

	requestProcessor IRequestProcessor[T]

	channel chan []IRequest[T]

	results chan CommandQueueResult[T]

	requestIds []uint64

	requestsProcessed []CommandQueueResult[T]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex
}

func NewCommandQueue[T schema.FixedSizeIdentifiable](numberOfWorkers int, batchSize int, bufferSize int, requestProcessor IRequestProcessor[T]) *CommandQueue[T] {

	return &CommandQueue[T]{

		numberOfWorkers:   numberOfWorkers,
		batchSize:         batchSize,
		bufferSize:        bufferSize,
		requestProcessor:  requestProcessor,
		channel:           make(chan []IRequest[T], bufferSize),
		results:           make(chan CommandQueueResult[T], bufferSize),
		requestIds:        make([]uint64, 0),
		requestsProcessed: make([]CommandQueueResult[T], 0),
	}
}

func (cq *CommandQueue[T]) NewRequestId() uint64 {

	return 0

	// id := lib.NewId(cq.requestIds)

	// cq.mu.Lock()

	// cq.requestIds = append(cq.requestIds, id)

	// cq.mu.Unlock()

	// return id
}

func (cq *CommandQueue[T]) Processor() IRequestProcessor[T] {

	return cq.requestProcessor
}

func (cq *CommandQueue[T]) StartWorkers() {

	for i := 0; i < cq.numberOfWorkers; i++ {

		cq.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer cq.requestsWaitGroup.Done()

			for batch := range cq.channel {

				for _, request := range batch {

					result := cq.requestProcessor.Process(request)

					cq.results <- CommandQueueResult[T]{request, result}
				}
			}
		}(i)
	}
	cq.resultsWaitGroup.Add(1)

	go func() {

		defer cq.resultsWaitGroup.Done()

		for result := range cq.results {

			// fmt.Println("%v", result)

			cq.mu.Lock()

			cq.requestsProcessed = append(cq.requestsProcessed, result)

			cq.mu.Unlock()
		}
	}()
}

func (cq *CommandQueue[T]) AddTask(task IRequest[T]) {

	batch := []IRequest[T]{task}

	cq.channel <- batch
}

func (cq *CommandQueue[T]) AddTasks(tasks []IRequest[T]) {

	numberOfBatches := (len(tasks) + cq.batchSize - 1) / cq.batchSize

	for i := 0; i < numberOfBatches; i++ {

		start := i * cq.batchSize

		end := (i + 1) * cq.batchSize

		if end > len(tasks) {

			end = len(tasks)
		}
		batch := tasks[start:end]

		cq.channel <- batch
	}
}

// func (cq *CommandQueue[T]) GetResults() []IResult[T] {

// 	var mu sync.Mutex

// 	var wg sync.WaitGroup

// 	results := []IResult[T]{}

// 	wg.Add(1)

// 	go func() {

// 		defer wg.Done()

// 		for result := range cq.results {

// 			mu.Lock()

// 			results = append(results, result)

// 			mu.Unlock()
// 		}
// 	}()

// 	close(cq.channel)

// 	cq.wg.Wait()

// 	close(cq.results)

// 	wg.Wait() // Wait for the results collection goroutine to finish

// 	return results
// }

func (cq *CommandQueue[T]) Stop() []CommandQueueResult[T] {

	close(cq.channel)

	cq.requestsWaitGroup.Wait()

	close(cq.results)

	cq.resultsWaitGroup.Wait()

	return cq.requestsProcessed
}
