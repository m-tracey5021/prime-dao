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

	requestIds []uint64

	// dependentProcesses map[uint64][]uint64 // e.g. process 0 depends on processes 1, 2, 3

	requestsProcessed []TSFQueueResult[T]

	requestsWaitGroup sync.WaitGroup

	resultsWaitGroup sync.WaitGroup

	mu sync.Mutex

	ctxmu sync.Mutex
}

func NewQueue[T any](numberOfWorkers int, batchSize int, bufferSize int) *TSFQueue[T] {

	return &TSFQueue[T]{

		numberOfWorkers: numberOfWorkers,

		batchSize: batchSize,

		bufferSize: bufferSize,

		channel: make(chan []IProcessableRequest[T], bufferSize),

		results: make(chan TSFQueueResult[T], bufferSize),

		requestIds: make([]uint64, 0), // do i need to initialise these?

		// dependentProcesses: make(map[uint64][]uint64, 0),

		requestsProcessed: make([]TSFQueueResult[T], 0),
	}
}

// func (queue *TSFQueue[T]) AddDependency(requestId uint64, dependencies []uint64) {

// 	_, found := queue.dependentProcesses[requestId]

// 	if found {

// 		queue.dependentProcesses[requestId] = append(queue.dependentProcesses[requestId], dependencies...)

// 	} else {

// 		queue.dependentProcesses[requestId] = dependencies
// 	}
// }

func (queue *TSFQueue[T]) Start(context *QueueContext[T]) {

	for i := 0; i < queue.numberOfWorkers; i++ {

		queue.requestsWaitGroup.Add(1)

		go func(workerId int) {

			defer queue.requestsWaitGroup.Done()

			for batch := range queue.channel {

				for _, request := range batch {

					// make sure dependencies are in different batches

					// so that if one comes before the other we dont have a race condition and block

					var processWaitGroup sync.WaitGroup

					processWaitGroup.Add(1)

					result := request.Process(&processWaitGroup, context)

					mappedResult := TSFQueueResult[T]{request.RequestId(), result}

					requestContext := RequestContext{request.RequestId(), &processWaitGroup}

					// queue.ctxmu.Lock()

					// context.completions = append(context.completions, signal)

					// queue.ctxmu.Unlock()

					context.completed <- requestContext

					queue.results <- mappedResult
				}
			}
		}(i)
	}
	// go func() {

	// 	for signal := context.completed {

	// 		// separate lock for this

	// 		context.completions = append(context.completions, signal)
	// 	}
	// }

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

	var dummyWaitGroup sync.WaitGroup

	dummyContext := QueueContext[T]{}

	return request.Process(&dummyWaitGroup, &dummyContext)
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
