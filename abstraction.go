package workerpool

import (
	"context"
)

// data serialization: using byte as generic data
// for go channel we can use encoding/gob to encode
// way for logging: logging can be handle by worker it self we can use goabs-log package for logging
// @TODO: way for error handling
// @TODO: testing
// @TODO: performance:w
// @TODO: check race condition

// Worker a function that receive context and data in []byte from queue then
// return an error if there are any unexpected problem happens the returned error
// can be treat in the different ways which based on the queue implementation.
// worker is one way processing (cannot return processed value), but you can delegate another message to
// passing outcome data to another queue to make a sequence of tasks.
type Worker = func(ctx context.Context, data []byte) error

// IDelegator interface of a delegator which delegate a queue item
// the queue item will be picked up by worker and processed
type IDelegator interface {
	// Delegate delegate a message to a queue which has
	// name match with the "queueName" argument
	Delegate(ctx context.Context, queueName string, data []byte) error
}

// IWorkerPool interface of a worker-pool which manage many worker
type IWorkerPool interface {
	// AddWorker add worker to handle a message of queue which has
	// name match with the "queueName" argument
	AddWorker(ctx context.Context, queueName string, worker Worker) error
}

// IQueue interface of a queue
type IQueue interface {
	// Enqueue add message to queue
	Enqueue(ctx context.Context, data []byte) error

	// Dequeue get message from queue. This function return 3 values
	// 1, `message` a read-only channel which will be feed by `Enqueue` function
	// 2, `error` a write-only channel which is used return error of worker function.
	// 3. `error` not nil if there are any problem when dequeue the message
	Dequeue(ctx context.Context) (<-chan []byte, chan<- error, error)
	
	// Dispose which will be used to clear up the queue when not using
	Dispose()
}
