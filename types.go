package pubsub

import (
	"time"
)

// Content ...
type Content struct {
	ExchangeName string
	QueueName    string
	Body         []byte
	Delay        time.Duration
}

// Consumer ...
type Consumer struct {
	ExchangeName string
	QueueName    string
	AutoAck      bool
	Receiver     chan []byte
}
