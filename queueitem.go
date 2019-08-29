package workerpool

import (
	"context"
)
// QueueItem implement IQueueItem which will be returned by Dequeue function
// Use NewQueueItem to build a QueueItem
type QueueItem struct{
	data []byte
	complete func(ctx context.Context, err error)
}

// NewQueueItem construct a queue item with buffer
// use WhenComplete chain function to hook the function which will be called
// when queue item was processed
func NewQueueItem(data []byte) *QueueItem {
	return &QueueItem{data: data}
}

// WhenComplete use to set hook function which will be called when queue item was processed
func (t *QueueItem) WhenComplete(complete func(ctx context.Context, err error)) *QueueItem {
	if t != nil {
		t.complete = complete
	}
	return t
}

// Data implement interface IQueueItem. Returns internal buffer of queue item
func (t *QueueItem) Data() []byte{
	return t.data
}

// Complete implement interface IQueueItem.
// this function will be call by queue when worker processed a queue item
func (t *QueueItem) Complete(ctx context.Context, err error) {
	if t.complete != nil {
		t.complete(ctx, err)
	}
}
