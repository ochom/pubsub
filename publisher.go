package pubsub

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher ...
type Publisher struct {
	url      string
	exchange string
	queue    string
	delay    time.Duration
}

// NewPublisher ...
func NewPublisher(rabbitURL, exchangeName, queueName string) *Publisher {
	return &Publisher{rabbitURL, exchangeName, queueName, 0}
}

// NewPublisherWithDelay ...
func NewPublisherWithDelay(rabbitURL, exchangeName, queueName string, delay time.Duration) *Publisher {
	return &Publisher{rabbitURL, exchangeName, queueName, delay}
}

// Publish ...
func (p *Publisher) Publish(body []byte) error {
	conn, ch, err := initQ(p.url)
	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, p.exchange, p.queue); err != nil {
		return err
	}

	headers := map[string]any{
		"x-delay": p.delay.Milliseconds(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish message to exchange
	err = ch.PublishWithContext(ctx,
		p.exchange, // exchange
		p.queue,    // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
		},
	)

	return err
}
