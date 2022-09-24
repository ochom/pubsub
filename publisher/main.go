package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ochom/pubsub"
	"github.com/streadway/amqp"
)

// Here we set the way error messages are displayed in the terminal.
func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func send(ch *amqp.Channel, body []byte, delay int64) error {

	return nil
}

func main() {
	// Let's catch the message from the terminal.
	reader := bufio.NewReader(os.Stdin)
	for {
		conf := pubsub.DefaultConfig()
		r, err := pubsub.NewRabbit(conf)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer r.Shutdown()

		fmt.Println("What message do you want to send?")
		mPayload, _ := reader.ReadString('\n')
		data := fmt.Sprintf(`{"payload": "%s"}`, strings.ReplaceAll(mPayload, "\n", ""))

		err = r.PublishWithDelay([]byte(data), 5000)
		failOnError(err, "Failed to publish a message")

		log.Printf(" [x] Congrats, sending message: %s", mPayload)
	}
}
