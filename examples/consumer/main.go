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

	consumer := pubsub.Consumer{
		ExchangeName: "test-exchange",
		QueueName:    "test-queue",
		AutoAck:      true,
		Messages:     make(chan []byte),
		Exit:         make(chan os.Signal, 1),
	}

	client = client.WithConsumer(consumer)

	go client.Consume()

	// handle messages
	go func() {
		for msg := range consumer.Messages {
			if err := processMessage(msg); err != nil {
				log.Println("Error: ", err.Error())
			}
		}
	}()

	log.Println("[*] Waiting for messages. To exit press CTRL+C")
	signal.Notify(consumer.Exit)
	<-consumer.Exit
}
