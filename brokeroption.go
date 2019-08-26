package workerpool

import (
	"github.com/jkaveri/goabs-workerpool/internal/abs"
)

// BrokerOption broker option is function that help to configure the broker
type BrokerOption = func (broker *Broker) error

// WithQueue register a queue with optional option
// you can bind worker with this queue by using WithWorker function
func WithQueue(name string, queue abs.IQueue, options ...QueueOption) BrokerOption {
	return func(broker *Broker) error {
		err := broker.qm.add(name, queue)
		if err != nil {
			return err
		}

		for _, opt := range options {
			err = opt(name, broker)
			if err != nil {
				return err
			}
		}

		return nil
	}
}
