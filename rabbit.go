package pubsub

import (
	"fmt"
	"time"

	"log"

	"github.com/streadway/amqp"
)

// Configuration ...
type Configuration struct {
	URL       string
	QueueName string
	Exchange  string
}

// DefaultConfig ...
func DefaultConfig() *Configuration {
	return &Configuration{
		URL:       "amqp://guest:guest@localhost:5672/",
		QueueName: "golang-queue",
		Exchange:  "golang-exchange",
	}
}

// NewConfig ...
func NewConfig(url, exchange, queueName string) *Configuration {
	return &Configuration{
		URL:       url,
		QueueName: queueName,
		Exchange:  exchange,
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

	return conn, ch, nil
}

// Rabbit ...
type Rabbit struct {
	conn       *amqp.Connection
	ch         *amqp.Channel
	exchange   string
	queueName  string
	closeChann chan *amqp.Error
	quitChann  chan bool
}

// NewRabbit ...
func NewRabbit(conf *Configuration) (*Rabbit, error) {
	conn, ch, err := initQ(conf.URL)
	if err != nil {
		return nil, err
	}

	r := Rabbit{
		conn:      conn,
		ch:        ch,
		exchange:  conf.Exchange,
		queueName: conf.QueueName,
	}

	err = r.initPubSub()
	if err != nil {
		return nil, err
	}

	r.quitChann = make(chan bool)
	r.closeChann = make(chan *amqp.Error)
	r.conn.NotifyClose(r.closeChann)

	go r.handleDisconnect()

	return &r, nil
}

func (r *Rabbit) initPubSub() error {
	// declare exchange
	args := make(amqp.Table)
	args["x-delayed-type"] = "direct"
	err := r.ch.ExchangeDeclare(
		r.exchange,          // name
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
	q, err := r.ch.QueueDeclare(
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
	err = r.ch.QueueBind(
		q.Name,     // queue name
		q.Name,     // routing key
		r.exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("queue Bind: %s", err.Error())
	}

	return nil
}

// handleDisconnect handle a disconnection trying to reconnect every 5 seconds
func (r *Rabbit) handleDisconnect() {
	for {
		select {
		case errChann := <-r.closeChann:
			if errChann != nil {
				log.Printf("error: rabbitMQ disconnection: %v", errChann)
			}
		case <-r.quitChann:
			r.conn.Close()
			log.Println("...rabbitMQ has been shut down")
			r.quitChann <- true
			return
		}

		log.Println("...trying to reconnect to rabbitMQ...")

		time.Sleep(5 * time.Second)

		if err := r.initPubSub(); err != nil {
			log.Printf("error: rabbitMQ error: %v", err)
		}
	}
}

// Shutdown closes rabbitmq's connection
func (r *Rabbit) Shutdown() {
	r.quitChann <- true

	log.Println("shutting down rabbitMQ's connection...")

	<-r.quitChann
}

// Publish ...
func (r *Rabbit) Publish(body []byte) error {
	return r.publish(body, 0)
}

// PublishWithDelay ...
func (r *Rabbit) PublishWithDelay(body []byte, delay int64) error {
	return r.publish(body, delay)
}

func (r *Rabbit) publish(body []byte, delay int64) error {
	headers := make(amqp.Table)
	if delay != 0 {
		headers["x-delay"] = delay
	}

	// publish message to exchange
	err := r.ch.Publish(
		r.exchange,  // exchange
		r.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Headers:      headers,
		},
	)

	return err
}

// Consume ...
func (r *Rabbit) Consume() (<-chan amqp.Delivery, error) {
	msgs, err := r.ch.Consume(
		r.queueName, // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	if err != nil {
		return nil, fmt.Errorf("queue Consume: %s", err.Error())
	}

	return msgs, nil
}
