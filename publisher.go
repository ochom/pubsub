package pubsub

import (
	"context"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publisher ...
type Publisher struct {
	url      string
	exchange string
	queue    string
}

// NewPublisher ...
func NewPublisher(rabbitURL, queueName string) *Publisher {
	exchange := fmt.Sprintf("%s-exchange", queueName)
	return &Publisher{rabbitURL, exchange, queueName}
}

// PublishWithDelay ...
func (p *Publisher) PublishWithDelay(body []byte, delay time.Duration) error {
	return p.publish(body, delay)
}

// Publish ...
func (p *Publisher) Publish(body []byte) error {
	return p.publish(body, 0)
}

// Publish ...
func (p *Publisher) publish(body []byte, delay time.Duration) error {
	conn, ch, err := initQ(p.url)
	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	if err := initPubSub(ch, p.exchange, p.queue); err != nil {
		return err
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
			Headers: amqp.Table{
				"x-delay": delay.Milliseconds(),
			},
		},
	)

	return err
}
