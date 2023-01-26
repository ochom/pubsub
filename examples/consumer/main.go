package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func processMessage(b []byte) error {
	log.Printf("received a message: %s", string(b))
	return nil
}

func main() {

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	client := pubsub.NewClient(rabbitURL)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	client = client.WithExit(exit)

	consumer := &pubsub.Consumer{
		ExchangeName: "test-exchange",
		QueueName:    "test-queue",
		AutoAck:      true,
		Receiver:     make(chan []byte),
	}

	consumer2 := &pubsub.Consumer{
		ExchangeName: "test-exchange",
		QueueName:    "test-queue2",
		AutoAck:      true,
		Receiver:     make(chan []byte),
	}

	consumers := []*pubsub.Consumer{consumer, consumer2}

	go client.Consume()

	// register consumers
	client.RegisterConsumers(consumers...)

	for _, c := range consumers {
		go func(c *pubsub.Consumer) {
			for msg := range c.Receiver {
				if err := processMessage(msg); err != nil {
					log.Println("Error: ", err.Error())
				}
			}
		}(c)
	}

	log.Println("[*] Waiting for messages. To exit press CTRL+C")
	<-exit
}
