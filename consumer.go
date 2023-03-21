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

// GetQueueName ...
func (c *Consumer) GetQueueName() string {
	return c.queue
}

// Consume consume messages from the channels
func (c *Consumer) Consume() (chan []byte, error) {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, c.exchange, c.queue); err != nil {
		return nil, fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
	}

	deliveries, err := ch.Consume(
		c.queue, // queue
		"",      // consumer
		true,    // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		return nil, fmt.Errorf("failed to consume messages: %s", err.Error())
	}

	// Create a new channel to send message bodies
	bodyCh := make(chan []byte)

	// Start a separate goroutine to send message bodies to the channel
	go func() {
		defer close(bodyCh)
		for d := range deliveries {
			bodyCh <- d.Body
		}
	}()

	// Return the channel to the caller
	return bodyCh, nil
}
