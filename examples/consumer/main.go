package main

import (
	"log"
	"os"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func processMessage(msg []byte) error {
	log.Printf("Received a message: %s", string(msg))
	return nil
}

func main() {

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	consumer := pubsub.NewConsumer(rabbitURL, "test-queue")

	msgs, err := consumer.Consume()
	if err != nil {
		log.Fatalf("failed to consume messages: %s", err.Error())
	}

	for i := 0; i < 5; i++ {
		go func(workerID int) {
			for msg := range msgs {
				if err := processMessage(msg); err != nil {
					log.Printf("failed to process a message: %s", err.Error())
				}
			}
		}(i)
	}

	log.Println("[*] Waiting for messages. To exit press CTRL+C")

	wait := make(chan os.Signal, 1)
	<-wait
}
