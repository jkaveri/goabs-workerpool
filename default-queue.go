package workerpool

import (
	"context"
	"fmt"
)

type Queue struct {
	data   chan []byte
	closed chan struct{}
}

var _ IQueue = (*Queue)(nil)

func NewQueue(size int) *Queue {
	return &Queue{data: make(chan []byte, size), closed: make(chan struct{})}
}

func (q Queue) Enqueue(ctx context.Context, data []byte) error {
	q.data <- data
	return nil
}

func (q Queue) Dequeue(ctx context.Context) (<-chan []byte, chan<- error, error) {
	errChan := make(chan error, 1)
	go func() {
		select {
		case err := <-errChan:
			fmt.Println("error", err)
		case <-ctx.Done():
			return
		case <-q.closed:
			return
		}
	}()
	return q.data, errChan, nil
}

func (q Queue) Dispose() {
	close(q.data)
	close(q.closed)
}
