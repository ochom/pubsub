package pubsub

import (
	"fmt"
	"log"
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
func (c *Consumer) Consume() (<-chan []byte, error) {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, c.exchange, c.queue); err != nil {
		return nil, fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
	}

	msgs, err := ch.Consume(
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

	deliveries := make(chan []byte)
	go func() {
		for msg := range msgs {
			delivery := msg.Body
			log.Println("Received a message: ", string(delivery))
			deliveries <- delivery
			log.Println("Done processing a message: ", string(delivery))
		}
	}()

	return deliveries, nil
}
