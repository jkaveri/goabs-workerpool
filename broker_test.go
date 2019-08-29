package workerpool

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goabs-workerpool/internal/abs"
	"github.com/jkaveri/goabs-workerpool/internal/mock"
)

func TestNewBroker(t *testing.T) {
	broker, err := NewBroker()

	assert.Nil(t, err)
	assert.NotNil(t, broker)
	assert.NotNil(t, broker.qm)
}

func TestNewBroker_WithOptions(t *testing.T) {
	testErr := errors.New("test_error")
	successOption := func(broker *Broker) error {
		return nil
	}
	errorOption := func(broker *Broker) error {
		return testErr
	}
	cases := []struct {
		name      string
		options   []BrokerOption
		err       error
		brokerNil bool
	}{
		{
			name: "success_option",
			options: []BrokerOption{
				successOption,
			},
			err:       nil,
			brokerNil: false,
		},
		{
			name: "multiple_success_options",
			options: []BrokerOption{
				successOption,
				successOption,
			},
			err:       nil,
			brokerNil: false,
		},
		{
			name: "error_option",
			options: []BrokerOption{
				errorOption,
			},
			err:       testErr,
			brokerNil: true,
		},
		{
			name: "error_and_success_option",
			options: []BrokerOption{
				errorOption,
				successOption,
			},
			err:       testErr,
			brokerNil: true,
		},
	}

	for i := 0; i < len(cases); i++ {
		item := cases[i]
		t.Run(item.name, func(t *testing.T) {
			broker, err := NewBroker(item.options...)
			assert.Equal(t, item.err, err)
			if item.brokerNil {
				assert.Nil(t, broker)
			} else {
				assert.NotNil(t, broker)
			}
		})
	}
}

func TestBroker_Assign(t *testing.T) {
	const qName = "test_queue"

	// build broker
	qm := newQueueMap()
	broker := &Broker{
		qm: qm,
	}
	_ = qm.add(qName, &fakeQueue{})

	// Action
	err := broker.Assign(context.TODO(), qName, fakeWorker, 1)

	// Assertion
	assert.Nil(t, err)
}

func TestBroker_AssignReturnError(t *testing.T) {
	// build worker
	qm := newQueueMap()
	broker := &Broker{
		qm: qm,
	}
	qName := "test_queue"

	// Action
	err := broker.Assign(context.TODO(), qName, fakeWorker, 1)

	// Assertion
	assert.Equal(t, ErrQueueNotExist, errors.Cause(err))
}

func TestBroker_AssignWhenDequeueError(t *testing.T) {
	// Arrange
	const qName = "test_queue"
	expectedErr := errors.New("some error")

	// mock queue
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	bgCtx := context.Background()
	ctx, _ := context.WithCancel(bgCtx)
	mockQueue := mock.NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(ctx).Return(nil, expectedErr)

	// build worker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{
		qm: qm,
	}

	// Action
	err := broker.Assign(bgCtx, qName, fakeWorker, 1)

	// Assertion
	assert.Equal(t, expectedErr, errors.Cause(err))
}

func TestBroker_Assign_ShouldBeCalled(t *testing.T) {
	// Arrange
	const qName = "test_queue"
	var wg sync.WaitGroup
	wg.Add(1) // add pending tasks

	qItemChan := make(chan abs.IQueueItem, 1)
	qItem := NewQueueItem([]byte("test"))

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock Queue
	mockQueue := mock.NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(gomock.Any()).Return(qItemChan, nil)

	// build worker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{
		qm: qm,
	}
	// mock worker
	worker := func(ctx context.Context, data []byte) error {
		assert.Equal(t, qItem.Data(), data) // assertion
		wg.Done()
		return nil
	}

	// Action
	err := broker.Assign(context.TODO(), qName, worker, 1)
	assert.Nil(t, err) // assertion
	// feed obj into obj-channel so the worker should receive the obj
	qItemChan <- qItem
	wg.Wait()
}

func TestBroker_Assign_CancelWorker(t *testing.T) {
	// Arrange
	const qName = "test_queue"

	qItemChan := make(chan IQueueItem, 1)
	qItem := NewQueueItem([]byte("test"))
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	called := false
	// Mock Queue
	mockQueue := mock.NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(gomock.Any()).Return(qItemChan, nil)

	// build worker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{
		qm: qm,
	}

	// mock worker
	worker := func(ctx context.Context, data []byte) error {
		called = true
		return nil
	}

	// Action
	err := broker.Assign(ctx, qName, worker, 1)
	assert.Nil(t, err) // assertion

	// cancel worker
	cancel()
	// makesure done channel to be selected before feed qItemChan
	time.Sleep(1 *time.Millisecond)
	qItemChan <- qItem
	assert.Equal(t, false, called)
	assert.Equal(t, 1, len(qItemChan))
	assert.Equal(t, 0, broker.Goroutines())
}

func TestBroker_Delegate(t *testing.T) {
	// arrange
	const qName = "test_queue"
	data := []byte("obj")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()

	// mock queue
	mockQueue := mock.NewMockQueue(ctrl)
	mockQueue.EXPECT().Enqueue(ctx, data).Return(nil)

	// build broker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{
		qm: qm,
	}

	// action
	err := broker.Delegate(ctx, qName, data)

	// assertion
	assert.Nil(t, err)
}

// this test catch the ErrorQueueNotExist when delegate message into
// queue which doesn't exist in the broker.
func TestBroker_Delegate_ReturnQueueNotExistError(t *testing.T) {
	// arrange
	const qName = "test_queue"

	// build broker
	qm := newQueueMap()
	broker := &Broker{
		qm: qm,
	}
	data := []byte("obj")

	// action
	err := broker.Delegate(context.TODO(), qName, data)

	// assertion
	assert.NotNil(t, err)
	assert.Equal(t, ErrQueueNotExist, errors.Cause(err))
}

func TestBroker_Cancel(t *testing.T) {
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

	broker := &Broker{
		qm:    qmap,
	}

	broker.Cancel(qname)

	assert.True(t, called)
}