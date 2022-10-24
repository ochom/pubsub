package pubsub

import "log"

// Consume consumes messages from the queue
// autoAck: true if the server should consider messages acknowledged once delivered; false if the server should expect explicit acknowledgements
func (r *Rabbit) Consume(consumer *Consumer) error {
	conn, ch, err := initQ(r.connectionURL)
	if err != nil {
		return err
	}
	defer ch.Close()
	defer conn.Close()

	err = r.initPubSub(ch)
	if err != nil {
		return err
	}

	msgs, err := ch.Consume(
		r.queueName,      // queue
		"",               // consumer
		consumer.AutoAck, // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	if err != nil {
		return err
	}

	go func() {
		for d := range msgs {
			if err := consumer.CallBack(consumer.Worker, d.Body); err == nil {
				if !consumer.AutoAck {
					_ = d.Ack(false)
				}
			}
		}
	}()

	log.Printf("Consumer: %d [*] Waiting for messages. To exit press CTRL+C", consumer.Worker)
	<-consumer.Exit
	return nil
}
