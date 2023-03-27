package pubsub

import (
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
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

	err = ch.Qos(1, 0, false) // fair dispatch
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
}

// initPubSub ...
func initPubSub(ch *amqp.Channel, args amqp.Table, channelType, exchangeName, queueName string) error {
	err := ch.ExchangeDeclare(
		exchangeName, // name
		channelType,  // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		args,         // arguments
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

// initDelayed ...
func initDelayed(ch *amqp.Channel, exchangeName, queueName string) error {
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	channelType := "x-delayed-message"
	return initPubSub(ch, args, channelType, exchangeName, queueName)
}

// initLazy ...
func initLazy(ch *amqp.Channel, exchangeName, queueName string) error {
	args := make(amqp.Table)
	args["x-queue-mode"] = "lazy"
	channelType := "direct"
	return initPubSub(ch, args, channelType, exchangeName, queueName)
}
