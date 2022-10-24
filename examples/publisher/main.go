package main

import (
	"fmt"
	"log"
	"math/rand"
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

	r := pubsub.NewRabbit(rabbitURL, exchangeName, queueName)

	for i := 0; i < 20; i++ {
		// random delay
		myDelay := rand.Intn(100) + delay

		data := fmt.Sprintf(`{"message": "%s", "job": %d, "delay": %d}`, message, i, myDelay)
		err := r.PublishWithDelay([]byte(data), int64(myDelay))
		if err != nil {
			log.Fatalf("Failed to publish a message: %s", err)
		}
	}

	log.Printf("All Message published")
}
