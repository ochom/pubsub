package main

import (
	"log"

	"github.com/ochom/pubsub"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conf := pubsub.DefaultConfig()
	r, err := pubsub.NewRabbit(conf)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer r.Shutdown()

	msgs, err := r.Consume()
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
