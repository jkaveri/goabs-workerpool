package workerpool

import (
	"context"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// QueueOption function configure queue
type QueueOption = func(name string, broker *Broker) error

// WithWorker shortcut of WithWorkers(ctx, worker, 1)
func WithWorker(ctx context.Context, worker abs.WorkerFunc, n int) QueueOption {
	return func(name string, broker *Broker) error {
		return broker.Assign(ctx, name, worker, n)
	}
}

