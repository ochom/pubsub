package pubsub

import (
	"fmt"
)

// Consumer ...
type Consumer struct {
	url      string
	exchange string
	queue    string
}

// NewConsumer ...
func NewConsumer(rabbitURL, queueName string) *Consumer {
	exchange := fmt.Sprintf("%s-exchange", queueName)
	return &Consumer{rabbitURL, exchange, queueName}
}

// Consume consume messages from the channels
func (c *Consumer) Consume(workerFunc func([]byte), isLazy bool) error {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}
	defer ch.Close()
	defer conn.Close()

	if isLazy {
		if err := initLazy(ch, c.exchange, c.queue); err != nil {
			return fmt.Errorf("failed to initialize a lazy pubsub: %s", err.Error())
		}
	} else {
		if err := initDelayed(ch, c.exchange, c.queue); err != nil {
			return fmt.Errorf("failed to initialize a instant pubsub: %s", err.Error())
		}
	}

	deliveries, err := ch.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		return fmt.Errorf("failed to consume messages: %s", err.Error())
	}

	for d := range deliveries {
		workerFunc(d.Body)
		d.Ack(false)
	}

	return nil
}
