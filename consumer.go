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
	workers  int
}

// NewConsumer ...
func NewConsumer(rabbitURL, queueName string, handler ConsumerHandler, workers int) *Consumer {
	exchange := fmt.Sprintf("%s-exchange", queueName)
	if workers <= 0 {
		workers = 1
	}
	return &Consumer{rabbitURL, exchange, queueName, handler, workers}
}

// GetQueueName ...
func (c *Consumer) GetQueueName() string {
	return c.queue
}

// GetWorkers ...
func (c *Consumer) GetWorkers() int {
	return c.workers
}

// Consume consume messages from the channels
func (c *Consumer) Consume(workerID int) error {
	conn, ch, err := initQ(c.url)
	if err != nil {
		return fmt.Errorf("failed to initialize a connection: %s", err.Error())
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, c.exchange, c.queue); err != nil {
		return fmt.Errorf("failed to initialize a pubsub: %s", err.Error())
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
		return fmt.Errorf("failed to consume messages: %s", err.Error())
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
