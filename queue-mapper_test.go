package workerpool

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueueMapper(t *testing.T) {
	mapper := newQueueMap()

	assert.NotNil(t, mapper)
}

func TestQueueMapper_Add(t *testing.T) {
	queueName := "process"
	mapper := &queueMap{
		queues: map[string]IQueue{},
	}
	q := &fakeQueue{}
	err := mapper.add(queueName, q)
	assert.Nil(t, err)
	assert.Equal(t, q, mapper.queues[queueName])
}
