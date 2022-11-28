package main

import (
	"fmt"
	event "listener-service/events"
	"log"
	"math"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	/* Try to connect to RabbitMQ */
	rabbitConn, err := connect()
	if err != nil {
		log.Fatal(err)
	}

	defer rabbitConn.Close()

	/* Start listening for messages */
	log.Println("Listening for and consuming RabbitMQ messages...")

	/* Create a consumer that consumes messages from the queue */
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	/* Watch the queue and consume events */
	err = consumer.Listen([]string{"log.INFO", "log.WARN", "log.ERROR"})

	if err != nil {
		log.Println(err)
	}
}

func connect() (*amqp.Connection, error) {
	var retries int64
	var connection *amqp.Connection
	backoff := time.Second * 1

	for {
		connAttempt, err := amqp.Dial("amqp://guest:guest@localhost")
		if err != nil {
			fmt.Println("RabbitMQ not ready")
			retries++
			if retries > 5 {
				fmt.Println(err)
				return nil, err
			}
			backoff = time.Duration(math.Pow(float64(backoff), 2)) * time.Second
			fmt.Println("Backing off...")
			time.Sleep(backoff)
		} else {
			connection = connAttempt
			break
		}
	}
	return connection, nil
}
