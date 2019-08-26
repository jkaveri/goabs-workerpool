package workerpool

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestWithQueue(t *testing.T) {
	const qName = "test"
	q := &fakeQueue{}
	broker := &Broker{qm: newQueueMap()}
	opt := WithQueue(qName, q)

	assert.NotNil(t, opt)

	err := opt(broker)

	assert.Nil(t, err)
	assert.Equal(t, broker.qm.queues[qName], q)
}

func TestWithQueue_MapperReturnErr(t *testing.T) {
	const qName = "test"
	q := &fakeQueue{}
	qm := newQueueMap()
	_ = qm.add(qName, q)
	broker := &Broker{qm: qm}
	opt := WithQueue(qName, q)

	assert.NotNil(t, opt)

	err := opt(broker)

	assert.Equal(t, ErrQueueAlreadyExist, errors.Cause(err))
}

func TestWithQueue_QueueOption(t *testing.T) {
	testErr := errors.New("test")
	testcases := []struct{
		name string
		queueOption func(called *bool) QueueOption
		err error
	} {
		{
			"should_be_called",
			func(called *bool) QueueOption {
				return func(name string, broker *Broker) error {
					*called = true
					return nil
				}
			},
			nil,
		},
		{
			"test_err",
			func(called *bool) QueueOption {
				return func(name string, broker *Broker) error {
					*called = true
					return testErr
				}
			},
			testErr,
		},
	}

	for i := 0; i < len(testcases); i++ {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			const qName = "test"
			q := &fakeQueue{}
			qm := newQueueMap()
			broker := &Broker{qm: qm}
			called := false
			opt := WithQueue(qName, q, tc.queueOption(&called))

			assert.NotNil(t, opt)

			err := opt(broker)

			assert.Equal(t, tc.err, err)
			assert.True(t, called)
		})
	}
}
