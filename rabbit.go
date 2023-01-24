package pubsub

import (
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Client ...
type Client struct {
	connectionURL string
	consumer      Consumer
}

// NewClient ...
func NewClient(url string) *Client {
	if url == "" {
		url = "amqp://guest:guest@localhost:5672/"
	}

	return &Client{
		connectionURL: url,
	}
}

func initQ(url string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	err = ch.Qos(1, 0, false) // fair dispatch
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

func (c *Client) initPubSub(ch *amqp.Channel, exchangeName, queueName string) error {
	// declare exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	err := ch.ExchangeDeclare(
		exchangeName,        // name
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
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("queue Declare: %s", err.Error())
	}

	// bind queue to exchange
	err = ch.QueueBind(
		q.Name,       // queue name
		q.Name,       // routing key
		exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
}

// WithConsumer ..
func (c *Client) WithConsumer(consumer Consumer) *Client {
	c.consumer = consumer
	return c
}

// Consume ...
func (c *Client) Consume() {
	conn, ch, err := initQ(c.connectionURL)
	if err != nil {
		log.Fatal(err)
	}

	defer ch.Close()
	defer conn.Close()

	err = c.initPubSub(ch, c.consumer.ExchangeName, c.consumer.QueueName)
	if err != nil {
		log.Fatal(err)
	}

	msgs, err := ch.Consume(
		c.consumer.QueueName, // queue
		"",                   // consumer
		c.consumer.AutoAck,   // auto-ack
		false,                // exclusive
		false,                // no-local
		false,                // no-wait
		nil,                  // args
	)

	if err != nil {
		log.Fatal(err)
	}

	// consume messages
	go func() {
		for d := range msgs {
			c.consumer.Messages <- d.Body
		}
	}()

	<-c.consumer.Exit
}
