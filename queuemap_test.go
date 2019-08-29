package workerpool

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goabs-workerpool/internal/mock"
)

func TestNewQueueMapper(t *testing.T) {
	mapper := newQueueMap()

	assert.NotNil(t, mapper)
}

func TestQueueMapper_Add(t *testing.T) {
	queueName := "process"
	mapper := &queueMap{
		queues: make(map[string]*queueMapItem),
	}
	q := &fakeQueue{}
	err := mapper.add(queueName, q)
	assert.Nil(t, err)
	assert.NotNil(t, mapper.queues[queueName])
	assert.Equal(t, q, mapper.queues[queueName].queue)
}

func TestQueueMapper_AddAlreadyExist(t *testing.T) {
	queueName := "test"
	mapper := &queueMap{
		queues: map[string]*queueMapItem{
			queueName: {
				queue:  &fakeQueue{},
				cancel: nil,
			},
		},
	}

	err := mapper.add(queueName, &fakeQueue{})
	assert.Equal(t, ErrQueueAlreadyExist, errors.Cause(err))
}

func TestQueueMap_Get(t *testing.T) {
	const qName = "test"
	expectedQueue := &fakeQueue{}
	mapper := &queueMap{
		queues: map[string]*queueMapItem{
			qName: {
				queue: expectedQueue,
			},
		},
	}

	item := mapper.get(qName)
	assert.NotNil(t, item)
	assert.Equal(t, expectedQueue, item.queue)
}

func TestQueueMap_Remove(t *testing.T) {
	const qName = "test"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueue := mock.NewMockQueue(ctrl)
	mockQueue.EXPECT().Dispose().Times(1)

	mapper := &queueMap{
		queues: map[string]*queueMapItem{
			qName: {
				queue:  mockQueue,
				cancel: nil,
			},
		},
	}

	mapper.remove(qName)

	assert.Nil(t, mapper.queues[qName])
}

func TestQueueMap_Cancel(t *testing.T) {
	const qname = "test"
	called := false
	qmap := &queueMap{queues: map[string]*queueMapItem{
		qname: {
			queue:  &fakeQueue{},
			cancel: func() {
				called = true
			},
		},
	}}

	qmap.cancel(qname)

	assert.True(t, called)
}

func TestQueueMap_CancelThenCreateContext(t *testing.T) {
	const qname = "test"
	called := false
	ctx := context.Background()
	qmap := &queueMap{queues: map[string]*queueMapItem{
		qname: {
			queue:  &fakeQueue{},
			cancel: func() {
				called = true
			},
		},
	}}

	newCtx := qmap.cancelThenCreateContext(ctx, qname)

	assert.True(t, called)
	assert.NotNil(t, newCtx.Done())
}
