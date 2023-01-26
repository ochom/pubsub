package main

import (
	"log"
	"os"

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

	client = client.WithConsumer(consumer)
	client = client.WithConsumer(consumer2)

	go client.Consume()

	// handle messages
	go func() {
		for msg := range consumer.Receiver {
			if err := processMessage(msg); err != nil {
				log.Println("Error: ", err.Error())
			}
		}
	}()

	// handle messages
	go func() {
		for msg := range consumer2.Receiver {
			if err := processMessage(msg); err != nil {
				log.Println("Error: ", err.Error())
			}
		}
	}()

	log.Println("[*] Waiting for messages. To exit press CTRL+C")
	<-exit
}
