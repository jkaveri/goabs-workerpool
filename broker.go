package workerpool

import (
	"context"

	"github.com/pkg/errors"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// Broker implement IWorkerPool and IDelegator
type Broker struct {
	qm *queueMap
}

var _ IDelegator = (*Broker)(nil)
var _ IWorkerPool = (*Broker)(nil)

// NewBroker create new broker with broker options
func NewBroker(opts ...BrokerOption) (*Broker, error) {
	b := &Broker{newQueueMap()}

	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

// AddWorker add worker to handle a message of queue which has
// name match with the "queueName" argument
func (t *Broker) AddWorker(ctx context.Context, qName string, worker abs.Worker) error {
	queue := t.qm.get(qName)
	if queue == nil {
		return errors.Wrap(ErrQueueNotExist, qName)
	}

	qData, err := queue.Dequeue(ctx)
	if err != nil {
		return errors.Wrap(err, "cannot dequeue when add worker")
	}

	go func() {
		for {
			select {
			case msg := <-qData:
				msg.OnComplete(ctx, worker(ctx, msg.GetData()))
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}

// Delegate delegate a message to a queue which has
// name match with the "queueName" argument
func (t *Broker) Delegate(ctx context.Context, name string, data []byte) error {
	queue := t.qm.get(name)
	if queue == nil {
		return errors.Wrap(ErrQueueNotExist, name)
	}
	return queue.Enqueue(ctx, NewMessage(data, nil))
}
