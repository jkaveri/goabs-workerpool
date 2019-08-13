package workerpool

type BrokerOption = func (broker *Broker) error

func WithQueue(name string, queue IQueue, options ...QueueOption) BrokerOption {
	return func(broker *Broker) error {
		err := broker.Add(name, queue)
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
