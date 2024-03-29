package workerpool

import (
	"context"
)

// MemoryQueue implements IQueue with use go-channel for communication between go-routines
type MemoryQueue struct {
	data   chan IQueueItem
	closed chan struct{}
}

// Compile time check where MemoryQueue implement IQueue or not
var _ IQueue = (*MemoryQueue)(nil)

// NewMemoryQueue create memory queue with channel buffer size
func NewMemoryQueue(size int) *MemoryQueue {
	return &MemoryQueue{
		data:   make(chan IQueueItem, size),
		closed: make(chan struct{}),
	}
}

// Enqueue add message to queue
func (q *MemoryQueue) Enqueue(ctx context.Context, data []byte) error {
	q.data <- NewQueueItem(data)
	return nil
}

// Dequeue get message from queue. This function return 3 values
// `message` a read-only channel which will be feed by `Enqueue` function
// `error` not nil if there are any problem when dequeue the message
func (q *MemoryQueue) Dequeue(ctx context.Context) (msgChan <-chan IQueueItem, err error) {
	return q.data, nil
}

// Dispose clean up resource by closing the channel
func (q *MemoryQueue) Dispose() {
	close(q.data)
	close(q.closed)
}
