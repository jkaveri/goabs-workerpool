package workerpool

import (
	"context"
)

// fakeQueue a fake queue for testing
type fakeQueue struct {
}

// Compile time check
var _ IQueue = (*fakeQueue)(nil)

// Enqueue add msg to queue
func (f *fakeQueue) Enqueue(ctx context.Context, data []byte) error {
	return nil
}

// Dequeue get msg from queue
// nolint gocritic
func (f *fakeQueue) Dequeue(ctx context.Context) (<-chan IQueueItem, error) {
	return nil, nil
}

// Dispose cleanup resource
func (f *fakeQueue) Dispose() {
}
