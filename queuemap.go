package workerpool

import (
	"sync"

	"github.com/pkg/errors"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// queueMap manage queue with race condition concern
type queueMap struct {
	queues map[string]abs.IQueue
	mutex  sync.Mutex
}

// newQueueMap create queue mapper
func newQueueMap() *queueMap {
	return &queueMap{
		queues: make(map[string]abs.IQueue),
	}
}

// add add queue with name
func (t *queueMap) add(name string, queue abs.IQueue) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if _, ok := t.queues[name]; ok {
		return errors.Wrap(ErrQueueAlreadyExist, name)
	}
	t.queues[name] = queue
	return nil
}

// remove queue by name
func (t *queueMap) remove(name string) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if queue, ok := t.queues[name]; ok {
		queue.Dispose()
		delete(t.queues, name)
	}
}

// get getting queue by name
func (t *queueMap) get(name string) abs.IQueue {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.queues[name]
}
