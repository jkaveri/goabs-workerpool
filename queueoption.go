package workerpool

import (
	"context"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// QueueOption function configure queue
type QueueOption = func(name string, broker *Broker) error

// WithWorkers create queue option which helps to add n workers to queue
func WithWorkers(ctx context.Context, worker abs.Worker, n int) QueueOption {
	if n <= 0 {
		n = 1
	}
	return func(name string, broker *Broker) error {
		for i := 0; i < n; i++ {
			err := broker.AddWorker(ctx, name, worker)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
// WithWorker shortcut of WithWorkers(ctx, worker, 1)
func WithWorker(ctx context.Context, worker abs.Worker) QueueOption {
	return WithWorkers(ctx, worker, 1)
}

