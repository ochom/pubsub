package pubsub

import (
	"fmt"

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
