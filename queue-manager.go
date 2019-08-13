package workerpool

import (
	"fmt"
	"sync"
)

type QueueMapper struct {
	queues map[string]IQueue
	mutex  sync.Mutex
}

func NewQueueMapper() *QueueMapper {
	return &QueueMapper{
		queues: make(map[string]IQueue),
	}
}

func (t *QueueMapper) Add(name string, queue IQueue) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	if _, ok := t.queues[name]; ok {
		return fmt.Errorf("queue %s already exist", name)
	}
	t.queues[name] = queue
	return nil
}

func (t *QueueMapper) Remove(name string, queue IQueue) error {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	if queue, ok := t.queues[name]; ok {
		queue.Dispose()
		delete(t.queues, name)
	}
	return nil
}

func (t *QueueMapper) Get(name string) IQueue {
	t.mutex.Lock()
	defer t.mutex.Unlock()
	return t.queues[name]
}
