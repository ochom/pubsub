package main

import (
	"log"
	"os"

	"github.com/ochom/pubsub"
)

func processMessage(worker int, b []byte) error {
	log.Printf("Worker %d: %s", worker, string(b))
	return nil
}

func main() {

	rabbitURL := os.Getenv("RABBIT_URL")
	exchangeName := os.Getenv("RABBIT_EXCHANGE")
	queueName := os.Getenv("RABBIT_QUEUE")

	r := pubsub.NewRabbit(rabbitURL, exchangeName, queueName)

	exit := make(chan bool)

	workers := 5
	for i := 0; i < workers; i++ {
		go func(worker int) {
			consumer := &pubsub.Consumer{
				Worker:   worker,
				Exit:     exit,
				AutoAck:  true,
				CallBack: processMessage,
			}

			err := r.Consume(consumer)
			if err != nil {
				log.Fatalf("Failed to consume a message: %s", err)
			}
		}(i)
	}

	<-exit
}
