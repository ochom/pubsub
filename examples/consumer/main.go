package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func processMessage(id int, jobs chan []byte) {
	for b := range jobs {
		log.Printf("worker: %d received a message: %s", id, string(b))
	}
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
		Workers:      5,
		Handler:      processMessage,
	}

	consumer2 := &pubsub.Consumer{
		ExchangeName: "test-exchange",
		QueueName:    "test-queue2",
		AutoAck:      true,
		Workers:      5,
		Handler:      processMessage,
	}

	consumers := []*pubsub.Consumer{consumer, consumer2}

	go client.Consume()

	// register consumers
	client.RegisterConsumers(consumers...)

	log.Println("[*] Waiting for messages. To exit press CTRL+C")
	<-exit
}
