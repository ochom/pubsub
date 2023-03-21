package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func main() {

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	consumer := pubsub.NewConsumer(rabbitURL, "test-queue")

	bodyCh, err := consumer.Consume()
	if err != nil {
		log.Fatalf("failed to consume messages: %s", err.Error())
	}

	go func() {
		for msg := range bodyCh {
			log.Printf("Received a message: %s", string(msg))
			// Why is no message logged here?

			// If you add a sleep here, the message will be logged
			time.Sleep(1 * time.Second)

		}
	}()

	// Use a buffered channel with capacity 1
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt)

	log.Println("Waiting for messages. To exit press CTRL+C")

	// Wait for a message on the stopCh channel
	<-stopCh
}
