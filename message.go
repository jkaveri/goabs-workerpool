package workerpool

import (
	"context"
)

// Message represent for a queue message which will be processed by worker
type Message struct {
	data           []byte
	onCompleteFunc func(ctx context.Context, err error)
}

// GetData return data in bytes
func (t *Message) GetData() []byte {
	return t.data
}

// OnComplete function will be called when message processing completed
func (t *Message) OnComplete(ctx context.Context, err error) {
	if t.onCompleteFunc != nil {
		t.onCompleteFunc(ctx, err)
	}
}

// NewMessage create a queue message
func NewMessage(data []byte, onComplete func(ctx context.Context, err error)) *Message {
	return &Message{
		data:           data,
		onCompleteFunc: onComplete,
	}
}