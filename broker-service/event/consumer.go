package event

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

/* Recieves events from the queue */
type Consumer struct {
	conn      *amqp.Connection
	queuename string
}

/*
Creates and returns a new consumer by initializing it's channel and exchange.
*/
func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer{
		conn: conn,
	}

	channel, err := consumer.conn.Channel()
	if err != nil {
		log.Println("Failed to create channel")
		return Consumer{}, err
	}

	/* Creates an exchange on the server. Idempotent. */
	err = declareExchange(channel)
	if err != nil {
		log.Println("Failed to create exchange")
		return Consumer{}, err
	}

	return consumer, nil
}

/*
Opens up a channel for the connection to communicate
with the RabbitMQ server, sets up a queue to hold messages
for the consumer, and binds all of the provided topics
to the queue so that the consumer can read them.
*/
func (consumer *Consumer) Listen(topics []string) error {
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := declareRandomQueue(ch)

	if err != nil {
		return err
	}

	/* Bind our channel to each of the topics */
	for _, s := range topics {
		ch.QueueBind(
			q.Name,
			s,
			"logs_topic",
			false,
			nil,
		)

		if err != nil {
			return err
		}
	}

	/* Start delivering queued messages through the channel */
	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)

	if err != nil {
		return err
	}

	forever := make(chan bool)

	/* Process each of the messages returned from the queue */
	go func() {

		for d := range messages {
			var payload Payload

			_ = json.Unmarshal(d.Body, &payload)

			go handlePayload(payload)

		}
	}()

	fmt.Printf("Waiting for messages on [Exchange, Queue] [logs_topic, %s]", q.Name)
	<-forever

	return nil

}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func handlePayload(payload Payload) {
	switch payload.Name {
	case "log", "event":
		// Log whatever we get
		err := logEvent(payload)
		if err != nil {
			log.Println(err)
		}
	default:
		log.Printf("Message type %s not supported", payload.Name)
	}
}

func logEvent(entry Payload) error {
	jsonData, _ := json.MarshalIndent(entry, "", "\t")
	logServiceUrl := "http://logger-service/log"

	request, err := http.NewRequest("POST", logServiceUrl, bytes.NewBuffer(jsonData))

	if err != nil {
		return err
	}

	client := http.Client{}

	response, err := client.Do(request)

	if err != nil {
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusAccepted {
		return errors.New(fmt.Sprintf("tried to log, got %d status", response.StatusCode))
	}

	return nil

}
