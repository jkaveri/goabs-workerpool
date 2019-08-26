package workerpool

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewMemoryQueue(t *testing.T) {
	const expectedSize = 10
	q := NewMemoryQueue(expectedSize)

	assert.NotNil(t, q)
	assert.Equal(t, expectedSize, cap(q.data))
	assert.NotNil(t, q.closed)
}

func TestMemoryQueue_Enqueue(t *testing.T) {
	expectedData := NewMessage([]byte("test"),nil)
	q := &MemoryQueue{
		data:   make(chan IMessage, 10),
		closed: make(chan struct{}),
	}

	err := q.Enqueue(context.TODO(), expectedData)
	assert.Nil(t, err)
	assert.Equal(t, expectedData, <-q.data)
}

func TestMemoryQueue_Dequeue(t *testing.T) {
	expectedData := NewMessage([]byte("test"),nil)
	q := &MemoryQueue{
		data:   make(chan IMessage, 10),
		closed: make(chan struct{}),
	}
	q.data <- expectedData

	msgChan, err := q.Dequeue(context.TODO())

	assert.Nil(t, err)
	assert.NotNil(t, msgChan)
	assert.Equal(t, expectedData, <- msgChan)
	assert.True(t, true)
}
