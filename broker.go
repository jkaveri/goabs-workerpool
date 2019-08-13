package workerpool

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
)

type Broker struct {
	*QueueMapper
}

func NewBroker(opts ...BrokerOption) (*Broker, error) {
	b := &Broker{NewQueueMapper()}

	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

func (t *Broker) AddWorker(ctx context.Context, name string, worker Worker) error {
	queue := t.QueueMapper.Get(name)
	if queue == nil {
		return errors.New("queue not exists")
	}

	qData, errChan, err := queue.Dequeue(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot dequeue when add worker")
	}

	go func() {
		for {
			select {
			case item := <-qData:
				errChan <- worker(ctx, item)
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

func (t *Broker) Delegate(ctx context.Context, name string, data []byte) error {
	queue := t.QueueMapper.Get(name)
	if queue == nil {
		return fmt.Errorf("queue \"%s\" not exists", name)
	}
	return queue.Enqueue(ctx, data)
}
