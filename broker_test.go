package workerpool

import (
	"context"
	"sync"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"

	"github.com/stretchr/testify/assert"
)

func TestNewBroker(t *testing.T) {
	broker, err := NewBroker()

	assert.Nil(t, err)
	assert.NotNil(t, broker)
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

func TestBroker_AddWorker(t *testing.T) {
	const qName = "test_queue"

	// build broker
	qm := newQueueMap()
	broker := &Broker{qm: qm}
	_ = qm.add(qName, &fakeQueue{})

	// Action
	err := broker.AddWorker(context.TODO(), qName, fakeWorker)

	// Assertion
	assert.Nil(t, err)
}

func TestBroker_AddWorkerReturnError(t *testing.T) {
	// build worker
	qm := newQueueMap()
	broker := &Broker{qm: qm}
	qName := "test_queue"

	// Action
	err := broker.AddWorker(context.TODO(), qName, fakeWorker)

	// Assertion
	assert.Equal(t, ErrQueueNotExist, errors.Cause(err))
}

func TestBroker_AddWorkerWhenDequeueError(t *testing.T) {
	// Arrange
	const qName = "test_queue"
	expectedErr := errors.New("some error")

	// mock queue
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()
	mockQueue := NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(ctx).Return(nil, nil, expectedErr)

	// build worker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{qm: qm}

	// Action
	err := broker.AddWorker(ctx, qName, fakeWorker)

	// Assertion
	assert.Equal(t, expectedErr, errors.Cause(err))
}

func TestBroker_AddWorker_ShouldBeCalled(t *testing.T) {
	// Arrange
	const qName = "test_queue"
	var wg sync.WaitGroup
	wg.Add(1) // add pending tasks

	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)
	expectedData := []byte("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()

	// Mock Queue
	mockQueue := NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(ctx).Return(dataChan, errChan, nil)

	// build worker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{qm: qm}

	// mock worker
	worker := func(ctx context.Context, data []byte) error {
		assert.Equal(t, expectedData, data) // assertion
		wg.Done()
		return nil
	}

	// Action
	err := broker.AddWorker(context.TODO(), qName, worker)
	assert.Nil(t, err) // assertion
	// feed data into data-channel so the worker should receive the data
	dataChan <- expectedData
	wg.Wait()
}

func TestBroker_AddWorker_ErrChanShouldBeFeed(t *testing.T) {
	// Arrange
	const qName = "test_queue"
	expectedErr := errors.New("expected")
	var wg sync.WaitGroup
	wg.Add(2)
	dataChan := make(chan []byte, 1)
	errChan := make(chan error, 1)
	expectedData := []byte("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()

	// mock queue
	mockQueue := NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(ctx).Return(dataChan, errChan, nil)

	// build broker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{qm: qm}

	// select the error channel for assertion
	go func() {
		err :=<- errChan
		assert.Equal(t, expectedErr, err) // assertion
		wg.Done()
	}()

	// mock worker
	worker := func(ctx context.Context, data []byte) error {
		assert.Equal(t, expectedData, data)
		wg.Done()
		return expectedErr // return error should be feed into error channel
	}

	// Action
	err := broker.AddWorker(ctx, qName, worker)
	assert.Nil(t, err) // assertion
	dataChan <- expectedData // feed data into the channel

	wg.Wait()
}

func TestBroker_AddWorker_CancelWorker(t *testing.T) {
	// Arrange
	const qName = "test_queue"

	dataChan := make(chan []byte)
	errChan := make(chan error)
	expectedData := []byte("test")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	notCall := true
	// Mock Queue
	mockQueue := NewMockQueue(ctrl)
	mockQueue.EXPECT().Dequeue(ctx).Return(dataChan, errChan, nil)

	// build worker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{qm: qm}

	// mock worker
	worker := func(ctx context.Context, data []byte) error {
		notCall = false
		return nil
	}

	// Action
	err := broker.AddWorker(ctx, qName, worker)
	assert.Nil(t, err) // assertion

	// cancel worker
	cancel()
	go func() {
		<- dataChan
	}()
	// feed data into data-channel so the worker should receive the data
	dataChan <- expectedData
	assert.Equal(t, true, notCall)
}

func TestBroker_Delegate(t *testing.T) {
	// arrange
	const qName = "test_queue"
	data := []byte("data")
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	ctx := context.TODO()

	// mock queue
	mockQueue := NewMockQueue(ctrl)
	mockQueue.EXPECT().Enqueue(ctx, data).Return(nil)

	// build broker
	qm := newQueueMap()
	_ = qm.add(qName, mockQueue)
	broker := &Broker{qm:qm}

	// action
	err := broker.Delegate(context.TODO(), qName, data)

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
	broker := &Broker{qm:qm}
	data := []byte("data")

	// action
	err := broker.Delegate(context.TODO(), qName, data)

	// assertion
	assert.NotNil(t, err)
	assert.Equal(t, ErrQueueNotExist, errors.Cause(err))
}


