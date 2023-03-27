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
	return p.publish(body, delay, false)
}

// Publish ...
func (p *Publisher) Publish(body []byte) error {
	return p.publish(body, 0, true)
}

// Publish ...
func (p *Publisher) publish(body []byte, delay time.Duration, isLazy bool) error {
	conn, ch, err := initQ(p.url)
	if err != nil {
		return err
	}

	defer ch.Close()
	defer conn.Close()

	headers := map[string]any{}
	if isLazy {
		if err := initLazy(ch, p.exchange, p.queue); err != nil {
			return err
		}
	} else {
		if err := initDelayed(ch, p.exchange, p.queue); err != nil {
			return err
		}
		headers["x-delay"] = delay.Milliseconds()
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
