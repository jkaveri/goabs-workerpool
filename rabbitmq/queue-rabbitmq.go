package rabbitmq

import (
	"context"

	"github.com/pkg/errors"
	"github.com/streadway/amqp"

	"github.com/jkaveri/goabs-workerpool"
)

// RabbitMQ rabbit-mq mediator
type RabbitMQ struct {
	name          string
	url           string
	channel       *amqp.Channel
	prefectCount  int
	ignoreFailure bool
}

var _ workerpool.IQueue = (*RabbitMQ)(nil)


func (t *RabbitMQ) Enqueue(ctx context.Context, msg workerpool.IMessage) error {

	err := t.channel.Publish(
		"",
		t.name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "plain/text",
			Body:         msg.GetData(),
		},
	)

	if err != nil {
		return errors.WithMessage(err, "Cannot publish message to queue")
	}

	return nil
}

func (t *RabbitMQ) Dequeue(ctx context.Context) (<-chan workerpool.IMessage, error) {
	delivery, err := t.channel.Consume(
		t.name,
		"",    // consumer name
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)

	if err != nil {
		return nil, errors.WithMessage(err, "Cannot consume queue")
	}
	result := make(chan workerpool.IMessage, 1)
	go func() {
		for {
			select {
			case d := <-delivery:
				result <- workerpool.NewMessage(d.Body, func(ctx context.Context, err error) {
					if err != nil {
						_ = d.Ack(false)
					} else {
						_ = d.Nack(false, false)
					}
				})
			case <-ctx.Done():
				return
			}
		}
	}()
	return result, nil
}

func (t *RabbitMQ) Dispose() {
	_ = t.channel.Close()
}

func (t *RabbitMQ) HandleWorkerResult(ct context.Context, err error) {
}


func (t *RabbitMQ) declareQueue(name string) (amqp.Queue, error) {
	return t.channel.QueueDeclare(
		name,
		true,
		false,
		false,
		false,
		nil,
	)
}
