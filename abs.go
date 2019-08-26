package workerpool

import (
	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// Worker a function that receive context and data in []byte from queue then
// return an error if there are any unexpected problem happens the returned error
// can be treat in the different ways which based on the queue implementation.
// worker is one way processing (cannot return processed value), but you can delegate another message to
// passing outcome data to another queue to make a sequence of tasks.
type Worker = abs.Worker

// IDelegator interface of a delegator which delegate a queue item
// the queue item will be picked up by worker and processed
type IDelegator = abs.IDelegator

// IWorkerPool interface of a worker-pool which manage many worker
type IWorkerPool = abs.IWorkerPool

// IQueue interface of a queue
type IQueue = abs.IQueue

// IMessage represent for a queue message
type IMessage = abs.IMessage