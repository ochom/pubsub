package pubsub

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish publish a message that will be consumed immediately
func (r *Rabbit) Publish(body []byte) error {
	return r.publish(body, 0)
}

// PublishWithDelay publish a message with delay in seconds
func (r *Rabbit) PublishWithDelay(body []byte, delay int64) error {
	return r.publish(body, delay)
}

func (r *Rabbit) publish(body []byte, delay int64) error {
	conn, ch, err := initQ(r.connectionURL)
	if err != nil {
		return err
	}
	defer ch.Close()
	defer conn.Close()

	err = r.initPubSub(ch)
	if err != nil {
		return err
	}

	headers := map[string]any{}
	if delay != 0 {
		headers["x-delay"] = delay * 1000 // convert to milliseconds
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// publish message to exchange
	err = ch.PublishWithContext(ctx,
		r.exchangeName, // exchange
		r.queueName,    // routing key
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
		},
	)
	return err
}
