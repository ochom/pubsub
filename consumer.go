package pubsub

import (
	"fmt"
	"log"
)

// ConsumerHandler is a function that handles the jobs
type ConsumerHandler func(msg []byte) error

// Consumer ...
type Consumer struct {
	url      string
	exchange string
	queue    string
	handler  ConsumerHandler
}

// NewConsumer ...
func NewConsumer(rabbitURL, exchangeName, queueName string, handler ConsumerHandler) *Consumer {
	return &Consumer{rabbitURL, exchangeName, queueName, handler}
}

// Consume consume messages from the channels
func (c *Consumer) Consume() error {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, c.exchange, c.queue); err != nil {
		return err
	}

	msgs, err := ch.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)

	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err.Error())
	}

	// consume messages
	for msg := range msgs {
		if err := c.handler(msg.Body); err != nil {
			// re-queue the message
			if err := msg.Nack(false, true); err != nil {
				log.Println("failed to re-queue a message")
				continue
			}

			log.Println("message re-queued")
			continue
		}

		// ack the message
		if err := msg.Ack(false); err != nil {
			log.Println("failed to ack a message")
			continue
		}

	}

	return nil
}
