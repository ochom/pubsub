package pubsub

// Consumer ...
type Consumer struct {
	ExchangeName string
	QueueName    string
	Worker       int
	Exit         chan bool
	AutoAck      bool
	CallBack     func(int, []byte) error
}

// Content ...
type Content struct {
	ExchangeName string
	QueueName    string
	Body         []byte
	Delay        int64
}
