package workerpool

import (
	"sync"

	"github.com/pkg/errors"
)

// queueMap manage queue with race condition concern
type queueMap struct {
	queues map[string]IQueue
	mutex  sync.Mutex
}

// newQueueMap create queue mapper
func newQueueMap() *queueMap {
	return &queueMap{
		queues: make(map[string]IQueue),
	}
}

// add add queue with name
func (t *queueMap) add(name string, queue IQueue) error {
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
func (t *queueMap) get(name string) IQueue {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.queues[name]
}
