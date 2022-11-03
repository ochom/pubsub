package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ochom/pubsub"
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

	rabbitURL := os.Getenv("RABBIT_URL")
	exchangeName := os.Getenv("RABBIT_EXCHANGE")
	queueName := os.Getenv("RABBIT_QUEUE")

	r := pubsub.NewRabbit(rabbitURL)

	for i := 0; i < 20; i++ {
		cnt := &pubsub.Content{
			ExchangeName: exchangeName,
			QueueName:    queueName,
			Body:         []byte(fmt.Sprintf("%s %d", message, delay)),
			Delay:        int64(delay),
		}

		err := r.Publish(cnt)
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
	}

	log.Printf("All Message published")
}
