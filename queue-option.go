package workerpool

import (
	"context"
)

type QueueOption = func(name string, broker *Broker) error

func WithWorkers(ctx context.Context, worker Worker, size int) QueueOption {
	if size <= 0 {
		size = 1
	}
	return func(name string, broker *Broker) error {
		for i := 0; i < size; i++ {
			err := broker.AddWorker(ctx, name, worker)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
func WithWorker(ctx context.Context, worker Worker) QueueOption {
	return WithWorkers(ctx, worker, 1)
}

