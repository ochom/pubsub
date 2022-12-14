package pubsub

import "time"

// Consumer ...
type Consumer struct {
	ExchangeName string
	QueueName    string
	Exit         chan bool
	AutoAck      bool
	CallBack     func([]byte) error
}

// Content ...
type Content struct {
	ExchangeName string
	QueueName    string
	Body         []byte
	Delay        time.Duration
}
