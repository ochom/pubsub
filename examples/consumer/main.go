package main

import (
	"log"
	"os"

	"github.com/ochom/pubsub"
)

func processMessage(b []byte) error {
	log.Printf("received a message: %s", string(b))
	return nil
}

func main() {

	rabbitURL := os.Getenv("RABBIT_URL")
	exchangeName := os.Getenv("RABBIT_EXCHANGE")
	queueName := os.Getenv("RABBIT_QUEUE")

	r := pubsub.NewRabbit(rabbitURL)

	exit := make(chan bool)

	workers := 5
	for i := 0; i < workers; i++ {
		go func(worker int) {
			consumer := &pubsub.Consumer{
				ExchangeName: exchangeName,
				QueueName:    queueName,
				Exit:         exit,
				AutoAck:      true,
				CallBack:     processMessage,
			}

			err := r.Consume(consumer)
			if err != nil {
				log.Fatalf("Failed to consume a message: %s", err)
			}
		}(i)
	}

	log.Printf("Consumers: %d [*] Waiting for messages. To exit press CTRL+C", workers)
	<-exit
}
