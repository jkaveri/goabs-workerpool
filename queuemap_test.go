package workerpool

import (
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
	"github.com/jkaveri/goabs-workerpool/internal/mock"
)

func TestNewQueueMapper(t *testing.T) {
	mapper := newQueueMap()

	assert.NotNil(t, mapper)
}

func TestQueueMapper_Add(t *testing.T) {
	queueName := "process"
	mapper := &queueMap{
		queues: map[string]abs.IQueue{},
	}
	q := &fakeQueue{}
	err := mapper.add(queueName, q)
	assert.Nil(t, err)
	assert.Equal(t, q, mapper.queues[queueName])
}

func TestQueueMapper_AddAlreadyExist(t *testing.T) {
	queueName := "test"
	mapper := &queueMap{
		queues: map[string]abs.IQueue{
			queueName: &fakeQueue{},
		},
		mutex:  sync.Mutex{},
	}

	err := mapper.add(queueName, &fakeQueue{})
	assert.Equal(t, ErrQueueAlreadyExist, errors.Cause(err))
}

func TestQueueMap_Get(t *testing.T) {
	const qName = "test"
	expectedQueue := &fakeQueue{}
	mapper := &queueMap{
		queues: map[string]abs.IQueue{
			qName: expectedQueue,
		},
		mutex:  sync.Mutex{},
	}

	queue := mapper.get(qName)
	assert.Equal(t, expectedQueue, queue)
}

func TestQueueMap_Remove(t *testing.T) {
	const qName = "test"
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockQueue := mock.NewMockQueue(ctrl)
	mockQueue.EXPECT().Dispose().Times(1)

	mapper := &queueMap{
		queues: map[string]abs.IQueue{
			qName: mockQueue,
		},
		mutex:  sync.Mutex{},
	}

	mapper.remove(qName)

	assert.Nil(t, mapper.queues[qName])
}
