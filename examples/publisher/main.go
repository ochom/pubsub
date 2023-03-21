package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ochom/pubsub"
	"github.com/ochom/pubsub/examples"
)

func fromArgs(args []string) (string, int) {
	if len(args) < 3 {
		log.Fatalf("Usage: %s <message> <delay>", args[0])
	}

	delay, err := strconv.Atoi(args[len(args)-1])
	if err != nil {
		log.Fatalf("Invalid delay: %s", args[2])
	}

	msg := strings.Join(args[1:len(args)-1], " ")

	return msg, delay
}

func main() {

	message, delay := fromArgs(os.Args)

	rabbitURL := examples.GetEnv("RABBIT_URL", "amqp://guest:guest@localhost:5672/")

	actualDelay := time.Duration(time.Second * time.Duration(delay))

	fmt.Printf("message will be published after %d ms\n", actualDelay.Milliseconds())

	p := pubsub.NewPublisher(rabbitURL, "test-queue")

	wg := sync.WaitGroup{}
	wg.Add(20)
	for i := 0; i < 20; i++ {
		go func() {
			defer wg.Done()
			if err := p.PublishWithDelay([]byte(message), actualDelay); err != nil {
				log.Fatalf("Failed to publish a message: %s", err)
			}
		}()
	}

	wg.Wait()
	log.Println("All Message published")
}
