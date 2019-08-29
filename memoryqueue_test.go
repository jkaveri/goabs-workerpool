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
	expectedData := []byte("test")
	q := &MemoryQueue{
		data:   make(chan IQueueItem, 10),
		closed: make(chan struct{}),
	}

	err := q.Enqueue(context.TODO(), expectedData)
	assert.Nil(t, err)
	assert.Equal(t, expectedData, (<-q.data).Data())
}

func TestMemoryQueue_Dequeue(t *testing.T) {
	expectedData := []byte("test")
	q := &MemoryQueue{
		data:   make(chan IQueueItem, 10),
		closed: make(chan struct{}),
	}
	q.data <- NewQueueItem(expectedData)

	msgChan, err := q.Dequeue(context.TODO())

	assert.Nil(t, err)
	assert.NotNil(t, msgChan)
	assert.Equal(t, expectedData, (<- msgChan).Data())
	assert.True(t, true)
}

func TestMemoryQueue_Dispose(t *testing.T) {
	q := &MemoryQueue{
		data:   make(chan IQueueItem, 10),
		closed: make(chan struct{}),
	}

	q.Dispose()

	_, ok := <- q.data
	assert.False(t, ok)

	_, ok = <- q.closed
	assert.False(t, ok)
}