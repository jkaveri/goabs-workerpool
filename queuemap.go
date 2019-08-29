package workerpool

import (
	"context"

	"github.com/pkg/errors"
)

type queueMapItem struct {
	queue  IQueue
	// cancel function will be use to cancel all worker that handle the queue
	cancel context.CancelFunc
}

// queueMap manage queue with race condition concern
type queueMap struct {
	queues map[string]*queueMapItem
}

// newQueueMap create queue mapper
func newQueueMap() *queueMap {
	return &queueMap{
		queues: make(map[string]*queueMapItem),
	}
}

// add add queue with name
func (t *queueMap) add(name string, queue IQueue) error {
	if _, ok := t.queues[name]; ok {
		return errors.Wrap(ErrQueueAlreadyExist, name)
	}
	t.queues[name] = &queueMapItem{
		queue:  queue,
		cancel: nil,
	}
	return nil
}

// remove queue by name
func (t *queueMap) remove(name string) {
	if item, ok := t.queues[name]; ok {
		if item.cancel != nil {
			item.cancel()
		}
		item.queue.Dispose()
		delete(t.queues, name)
	}
}

// get getting queue by name
func (t *queueMap) get(name string) *queueMapItem {
	return t.queues[name]
}

// cancel cancel all worker to handle queue
func (t *queueMap) cancel(qName string) {
	if item, ok := t.queues[qName]; ok && item.cancel != nil {
		item.cancel()
		item.cancel = nil
	}
}

// cancelThenCreateContext
func (t *queueMap) cancelThenCreateContext(ctx context.Context, qName string) context.Context {
	newCtx, cancel := context.WithCancel(ctx)

	if item, ok := t.queues[qName]; ok {
		if item.cancel != nil {
			item.cancel()
		}

		item.cancel = cancel
	}

	return newCtx
}
