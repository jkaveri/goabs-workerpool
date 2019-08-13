package workerpool

import (
	"context"
)

// data serialization: using byte as generic data
// for go channel we can use encoding/gob to encode
// way for logging: logging can be handle by worker it self we can use goabs-log package for logging
// @TODO: way for error handling
// @TODO: testing
// @TODO: performance
// @TODO: check race condition

type Worker = func(context context.Context, data []byte) error

type IDelegator interface {
	Delegate(ctx context.Context, name string, data []byte) error
}

type IWorkerPool interface {
	AddWorker(ctx context.Context, name string, worker Worker) error
}

type IQueue interface {
	Enqueue(ctx context.Context, data []byte) error
	Dequeue(ctx context.Context) (<-chan []byte, chan <- error, error)
	Dispose()
}
