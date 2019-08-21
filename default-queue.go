package workerpool

import (
	"context"
	"fmt"
)

// MemoryQueue implements IQueue with use go-channel for communication between go-routines
type MemoryQueue struct {
	data   chan []byte
	closed chan struct{}
}

// Compile time check where MemoryQueue implement IQueue or not
var _ IQueue = (*MemoryQueue)(nil)

// NewQueue create memory queue with channel buffer size
func NewQueue(size int) *MemoryQueue {
	return &MemoryQueue{data: make(chan []byte, size), closed: make(chan struct{})}
}

// Enqueue add message to queue
func (q MemoryQueue) Enqueue(ctx context.Context, data []byte) error {
	q.data <- data
	return nil
}

// Dequeue get message from queue. This function return 3 values
// 1, `message` a read-only channel which will be feed by `Enqueue` function
// 2, `error` a write-only channel which is used return error of worker function.
// 3. `error` not nil if there are any problem when dequeue the message
func (q MemoryQueue) Dequeue(ctx context.Context) (msgChan <-chan []byte, errChan chan<- error, err error) {
	tempErrChan := make(chan error, 1)
	go func() {
		select {
		case err := <-tempErrChan:
			fmt.Println("error", err)
		case <-ctx.Done():
			return
		case <-q.closed:
			close(tempErrChan)
			return
		}
	}()
	return q.data, tempErrChan, nil
}

// Dispose clean up resource by closing the channel
func (q MemoryQueue) Dispose() {
	close(q.data)
	close(q.closed)
}
