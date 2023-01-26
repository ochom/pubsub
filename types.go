package pubsub

import (
	"time"
)

// ConsumerName ...
type ConsumerName string

// Content ...
type Content struct {
	ExchangeName string
	QueueName    string
	Body         []byte
	Delay        time.Duration
}

// Consumer ...
type Consumer struct {
	Name         ConsumerName
	ExchangeName string
	QueueName    string
	AutoAck      bool
	Workers      int
	Handler      func(id int, msg chan []byte)
}
