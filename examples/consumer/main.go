package main

import (
	"log"
	"os"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func processMessage() pubsub.ConsumerHandler {
	return func(id int, jobs chan []byte) {
		for b := range jobs {
			log.Printf("worker: %d received a message: %s", id, string(b))
		}
	}
}

func main() {

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	consumer := pubsub.NewConsumer(rabbitURL, "test-exchange", "test-queue", true, 5, processMessage())

	go consumer.Consume()

	log.Println("[*] Waiting for messages. To exit press CTRL+C")

	wait := make(chan os.Signal, 1)
	<-wait
}
