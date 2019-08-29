package abs

import (
	"context"
)

// WorkerFunc is a function that receive and process queued task.
// error will be returned if any un-happy case happens.
// tip: WorkerFunc is one way processing (cannot return processed value),
// but you can delegate another message to passing outcome data to another
// queue to make a sequence of tasks.
type WorkerFunc = func(ctx context.Context, task []byte) error

// IDelegator delegator which delegate a task to a queue which will be handled by workers
type IDelegator interface {
	// Delegate delegates a task to a queue by "queueName"
	Delegate(ctx context.Context, queueName string, task []byte) error
}

// IWorkerPool interface of a worker-pool which manage many worker
type IWorkerPool interface {
	// Assign registers n workers to handle tasks of queue.
	// workerFunc which handle task in queue you can only register one workerFunc
	// for a queue. if there is a workerFunc registered already, it will be removed (not interrupt).
	// n is number of worker (goroutines) you want to create.
	Assign(ctx context.Context, qName string, workerFunc WorkerFunc, n int) error

	// Cancel un-schedule worker function of a queue if any.
	Cancel(qName string)
}

// IQueue a queue contains all tasks which has a same thing to do.
type IQueue interface {
	// Enqueue add task into queue
	Enqueue(ctx context.Context, data []byte) error

	// Dequeue get task from queue. This function returns
	// A channel will be returned which worker can use to pick a task
	// an error will be return if any un-happy case happens.
	Dequeue(ctx context.Context) (<-chan IQueueItem, error)

	// Dispose clean up un-managed resource
	Dispose()
}

// IQueueItem interface of queue item result of Dequeue method of a queue
type IQueueItem interface {
	// Data return binary data of queue item
	Data() []byte

	// Complete will be called when worker complete
	// err can be nil or not nil based on the worker result
	Complete(ctx context.Context, err error)
}
