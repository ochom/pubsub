package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func processMessage() pubsub.ConsumerHandler {
	return func(msg []byte) error {
		randomInt := randomInt(1, 10)
		if randomInt%2 == 0 {
			return fmt.Errorf("failed to process a message")
		}

		log.Printf("Received a message: %s", string(msg))
		return nil
	}
}

func main() {

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")
	consumer := pubsub.NewConsumer(rabbitURL, "test-queue", processMessage())

	workers := 5
	for i := 0; i < workers; i++ {
		go consumer.Consume(i)
	}

	log.Println("[*] Waiting for messages. To exit press CTRL+C")

	wait := make(chan os.Signal, 1)
	<-wait
}
