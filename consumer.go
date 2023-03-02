package pubsub

import "fmt"

// ConsumerHandler is a function that handles the jobs
type ConsumerHandler func(workerID int, jobs chan []byte)

// Consumer ...
type Consumer struct {
	url      string
	exchange string
	queue    string
	autoAck  bool
	workers  int
	handler  ConsumerHandler
}

// NewConsumer ...
func NewConsumer(rabbitURL, exchangeName, queueName string, autoAck bool, workers int, handler ConsumerHandler) *Consumer {
	return &Consumer{rabbitURL, exchangeName, queueName, autoAck, workers, handler}
}

// NewAutoAckConsumer ...
func NewAutoAckConsumer(rabbitURL, exchangeName, queueName string, workers int, handler ConsumerHandler) *Consumer {
	return &Consumer{rabbitURL, exchangeName, queueName, true, workers, handler}
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
		c.queue,   // queue
		"",        // consumer
		c.autoAck, // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)

	if err != nil {
		return fmt.Errorf("failed to register a consumer: %s", err.Error())
	}

	// create a receiver jobs channel
	jobs := make(chan []byte)

	for i := 0; i < c.workers; i++ {
		go c.handler(i, jobs)
	}

	// consume messages
	for msg := range msgs {
		jobs <- msg.Body
	}

	return nil
}
