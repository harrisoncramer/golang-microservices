package main

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "80"

type Config struct {
	Rabbit *amqp.Connection
}

func main() {

	rabbitConn, err := connect()

	if err != nil {
		log.Fatal(err)
	}

	/* The base receiver for our routes */
	app := Config{
		Rabbit: rabbitConn,
	}

	log.Printf("Starting broker service on %s", webPort)

	/* Define HTTP Server*/
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	/* Start the server */
	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

func connect() (*amqp.Connection, error) {
	var retries int64
	var connection *amqp.Connection
	backoff := time.Second * 1

	for {
		connAttempt, err := amqp.Dial("amqp://guest:guest@rabbitmq")
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
