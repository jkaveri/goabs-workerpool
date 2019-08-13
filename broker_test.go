package workerpool

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewBroker(t *testing.T) {
	broker, err := NewBroker()

	assert.Nil(t, err)
	assert.NotNil(t, broker)
}

func TestNewBroker_WithOptions(t *testing.T) {
	testErr := errors.New("test_error")
	var successOption BrokerOption = func(broker *Broker) error {
		return nil
	}
	var errorOption BrokerOption = func(broker *Broker) error {
		return testErr
	}
	cases := []struct{
		name string
		options []BrokerOption
		err error
		brokerNil bool
	}{
		{
			name: "success_option",
			options: []BrokerOption{
				successOption,
			},
			err: nil,
			brokerNil: false,
		},
		{
			name: "multiple_success_options",
			options: []BrokerOption{
				successOption,
				successOption,
			},
			err: nil,
			brokerNil: false,
		},
		{
			name: "error_option",
			options: []BrokerOption{
				errorOption,
			},
			err: testErr,
			brokerNil: true,
		},
		{
			name: "error_and_success_option",
			options: []BrokerOption{
				errorOption,
				successOption,
			},
			err: testErr,
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