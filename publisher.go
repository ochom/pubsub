package pubsub

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Publish publish a message to the channels
func (r *Rabbit) Publish(cnt *Content) error {
	conn, ch, err := initQ(r.connectionURL)
	if err != nil {
		return err
	}
	defer ch.Close()
	defer conn.Close()

	err = r.initPubSub(ch, cnt.ExchangeName, cnt.QueueName)
	if err != nil {
		return err
	}

	headers := map[string]any{}
	if cnt.Delay != 0 {
		headers["x-delay"] = cnt.Delay * 1000 // convert to milliseconds
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
