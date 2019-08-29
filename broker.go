package workerpool

import (
	"context"
	"sync"

	"github.com/pkg/errors"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// Broker implement IWorkerPool and IDelegator
type Broker struct {
	qm         *queueMap
	mutex      sync.Mutex
	goroutines int
}

var _ IDelegator = (*Broker)(nil)
var _ IWorkerPool = (*Broker)(nil)

// NewBroker create new broker
func NewBroker(opts ...BrokerOption) (*Broker, error) {
	b := &Broker{
		qm: newQueueMap(),
	}

	for _, opt := range opts {
		err := opt(b)
		if err != nil {
			return nil, err
		}
	}

	return b, nil
}

// Assign registers n workers to handle tasks of queue.
// workerFunc which handle task in queue you can only register one workerFunc
// for a queue. if there is a workerFunc registered already, it will be removed (not interrupt).
// n is number of worker (goroutines) you want to create.
// return error if a queue with name = qName not exist or any error happens when Dequeue.
func (t *Broker) Assign(ctx context.Context, qName string, worker abs.WorkerFunc, n int) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	item := t.qm.get(qName)
	if item == nil {
		return errors.Wrap(ErrQueueNotExist, qName)
	}

	// cancel exist queue and create new cancellation context
	ctx = t.qm.cancelThenCreateContext(ctx, qName)

	qItemChan, err := item.queue.Dequeue(ctx)
	if err != nil {
		// if there is an error cancel queue
		// some queue implementation can clean resource based on context cancellation
		t.qm.cancel(qName)
		return errors.Wrap(err, "cannot dequeue when add worker")
	}

	// create n worker
	for i := 0; i < n; i++ {
		// increase number of goroutines when create
		t.goroutines++

		go func() {
			for {
				select {
				case qItem := <-qItemChan:
					qItem.Complete(ctx, worker(ctx, qItem.Data()))
				case <-ctx.Done():
					t.mutex.Lock()
					t.goroutines--
					t.mutex.Unlock()
					return
				}
			}
		}()
	}

	return nil
}

// Cancel cancel scheduled workerFunc of a queue.
func (t *Broker) Cancel(qName string) {
	t.qm.cancel(qName)
}

// Delegate delegates a task to a queue by "queueName"
func (t *Broker) Delegate(ctx context.Context, name string, data []byte) error {
	item := t.qm.get(name)
	if item == nil {
		return errors.Wrap(ErrQueueNotExist, name)
	}
	return item.queue.Enqueue(ctx, data)
}

// Goroutines return number of goroutines were created
func (t *Broker) Goroutines() int {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.goroutines
}
