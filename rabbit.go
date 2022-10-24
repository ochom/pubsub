package pubsub

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func initQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func (r *Rabbit) initPubSub(ch *amqp.Channel) error {
	// declare exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	err := ch.ExchangeDeclare(
		r.exchangeName,      // name
		"x-delayed-message", // type
		true,                // durable
		false,               // auto-deleted
		false,               // internal
		false,               // no-wait
		args,                // arguments
	)
	if err != nil {
		return fmt.Errorf("exchange Declare: %s", err.Error())
	}

	// declare queue
	q, err := ch.QueueDeclare(
		r.queueName, // name
		true,        // durable
		false,       // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,         // queue name
		q.Name,         // routing key
		r.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
}

// Rabbit ...
type Rabbit struct {
	connectionURL string
	exchangeName  string
	queueName     string
}

// Consumer ...
type Consumer struct {
	Worker   int
	Exit     chan bool
	AutoAck  bool
	CallBack func(int, []byte) error
}

// NewRabbit ...
func NewRabbit(url, exchangeName, queueName string) *Rabbit {

	r := Rabbit{
		connectionURL: url,
		exchangeName:  exchangeName,
		queueName:     queueName,
	}

	return &r
}

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

	// publish message to exchange
	err = ch.Publish(
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

// Consume consumes messages from the queue
// autoAck: true if the server should consider messages acknowledged once delivered; false if the server should expect explicit acknowledgements
func (r *Rabbit) Consume(consumer *Consumer) error {
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

	msgs, err := ch.Consume(
		r.queueName,      // queue
		"",               // consumer
		consumer.AutoAck, // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := consumer.CallBack(consumer.Worker, d.Body); err == nil {
				if !consumer.AutoAck {
					d.Ack(false)
				}
			}
		}
	}()

	log.Printf("Consumer: %d [*] Waiting for messages. To exit press CTRL+C", consumer.Worker)
	<-consumer.Exit
	return nil
}
