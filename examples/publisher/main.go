package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func fromArgs(args []string) (string, int) {
	if len(args) != 3 {
		log.Fatalf("Usage: %s [message] [delay]", os.Args[0])
	}

	delay, err := strconv.Atoi(args[2])
	if err != nil {
		log.Fatalf("Invalid delay: %s", args[2])
	}

	return args[1], delay
}

func main() {

	message, delay := fromArgs(os.Args)

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")

	actualDelay := time.Duration(time.Second * time.Duration(delay))

	fmt.Printf("message will be published after %d ms\n", actualDelay.Milliseconds())

	client := pubsub.NewClient(rabbitURL)

	for i := 0; i < 20; i++ {
		cnt := &pubsub.Content{
			ExchangeName: "test-exchange",
			QueueName:    "test-queue",
			Body:         []byte(fmt.Sprintf("test-queue %s %d", message, i)),
			Delay:        actualDelay,
		}

		if err := client.Publish(cnt); err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
	}

	for i := 0; i < 20; i++ {
		cnt := &pubsub.Content{
			ExchangeName: "test-exchange",
			QueueName:    "test-queue2",
			Body:         []byte(fmt.Sprintf("test-queue2 %s %d", message, i)),
			Delay:        actualDelay,
		}

		if err := client.Publish(cnt); err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
	}

	log.Println("All Message published")
}
