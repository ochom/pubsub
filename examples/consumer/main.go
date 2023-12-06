package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/ochom/gutils/pubsub"
)

func main() {
	workerFunc := func(d []byte) {
		log.Printf("Received message: %s", string(d))
	}

	go func() {
		if err := pubsub.Consume("my-test-queue", workerFunc); err != nil {
			log.Fatalf("Failed to consume messages: %s", err.Error())
		}
	}()

	log.Println("Waiting for messages. To exit press CTRL+C")

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	// Wait for a message on the exit channel
	<-exit
}
