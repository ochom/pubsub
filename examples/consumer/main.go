package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ochom/pubsub"
)

func main() {

	rabbitURL := "amqp://admin:admin2020@localhost:5672/"
	consumer := pubsub.NewConsumer(rabbitURL, "test-queue")

	workerFunc := func(d []byte) {
		log.Printf("Received message: %s", string(d))
	}

	go func() {
		if err := consumer.Consume(workerFunc); err != nil {
			log.Fatalf("Failed to consume messages: %s", err.Error())
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	// Wait for a message on the exit channel
	<-exit
}
