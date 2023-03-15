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
func NewConsumer(rabbitURL, queueName string, handler ConsumerHandler) *Consumer {
	exchange := fmt.Sprintf("%s-exchange", queueName)
	return &Consumer{rabbitURL, exchange, queueName, handler}
}

// Consume consume messages from the channels
func (c *Consumer) Consume(workerID int) error {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return fmt.Errorf("[%s] worker [%d] failed to initialize a connection: %s", c.queue, workerID, err.Error())
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, c.exchange, c.queue); err != nil {
		return fmt.Errorf("[%s] worker [%d] failed to initialize a pubsub: %s", c.queue, workerID, err.Error())
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
		return fmt.Errorf("[%s] worker [%d] failed to consume messages: %s", c.queue, workerID, err.Error())
	}

	// consume messages
	for msg := range msgs {
		if err := c.handler(msg.Body); err != nil {
			// re-queue the message
			if err := msg.Nack(false, true); err != nil {
				log.Printf("[%s] worker [%d] failed to re-queue a message", c.queue, workerID)
				continue
			}

			log.Printf("[%s] worker [%d] re-queued a message", c.queue, workerID)
			continue
		}

		// ack the message
		if err := msg.Ack(false); err != nil {
			log.Printf("[%s] worker [%d] failed to ack a message", c.queue, workerID)
			continue
		}

	}

	return nil
}
