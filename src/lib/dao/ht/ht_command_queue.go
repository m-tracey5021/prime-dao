package ht

import (
	"sync"
	"transformer/src/lib"
	"transformer/src/lib/dao/schema"
)

type CommandQueue[T schema.FixedSizeIdentifiable] struct {
	numberOfWorkers int

	batchSize int

	channel chan []IRequest[T]

	results chan IResult[T]

	wg sync.WaitGroup

	requestsProcessed []uint64

	requestProcessor IRequestProcessor[T]
}

func NewCommandQueue[T schema.FixedSizeIdentifiable](numberOfWorkers int, batchSize int, requestProcessor IRequestProcessor[T]) *CommandQueue[T] {

	return &CommandQueue[T]{

		numberOfWorkers:   numberOfWorkers,
		batchSize:         batchSize,
		channel:           make(chan []IRequest[T], numberOfWorkers),
		results:           make(chan IResult[T], numberOfWorkers),
		requestsProcessed: make([]uint64, 0),
		requestProcessor:  requestProcessor,
	}
}

func (cq *CommandQueue[T]) Processor() IRequestProcessor[T] {

	return cq.requestProcessor
}

func (cq *CommandQueue[T]) NewRequestId() uint64 {

	return lib.NewId(cq.requestsProcessed)
}

func (cq *CommandQueue[T]) StartWorkers() {

	for i := 0; i < cq.numberOfWorkers; i++ {

		cq.wg.Add(1)

		go func(workerId int) {

			defer cq.wg.Done()

			for batch := range cq.channel {

				for _, request := range batch {

					result := cq.requestProcessor.Process(request)

					result.SetAssociatedRequestId(request.RequestId())

					cq.results <- result
				}
			}
		}(i)
	}
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

func (cq *CommandQueue[T]) GetResults() []IResult[T] {

	var mu sync.Mutex

	var wg sync.WaitGroup

	results := []IResult[T]{}

	wg.Add(1)

	go func() {

		defer wg.Done()

		for result := range cq.results {

			mu.Lock()

			results = append(results, result)

			mu.Unlock()
		}
	}()

	close(cq.channel)

	cq.wg.Wait()

	close(cq.results)

	wg.Wait() // Wait for the results collection goroutine to finish

	return results
}

func (cq *CommandQueue[T]) Stop() {

	close(cq.channel)

	cq.wg.Wait()

	close(cq.results)
}
