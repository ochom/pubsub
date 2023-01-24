package pubsub

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish publish a message to the channels
func (c *Client) Publish(cnt *Content) error {
	conn, ch, err := initQ(c.connectionURL)
	if err != nil {
		return err
	}
	defer ch.Close()
	defer conn.Close()

	err = c.initPubSub(ch, cnt.ExchangeName, cnt.QueueName)
	if err != nil {
		return err
	}

	headers := map[string]any{
		"x-delay": cnt.Delay.Milliseconds(),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish message to exchange
	err = ch.PublishWithContext(ctx,
		cnt.ExchangeName, // exchange
		cnt.QueueName,    // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         cnt.Body,
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
		},
	)
	return err
}
