package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

/* Pushes an event to RabbitMQ with the given severity */
func (e *Emitter) Push(event string, severity string) error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()

	err = channel.Publish(
		"logs_topic",
		severity,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(event),
		},
	)

	if err != nil {
		return err
	}

	return nil

}

/* Creates a new emitter with established channel and exchange with RabbitMQ */
func NewEventEmitter(conn *amqp.Connection) (Emitter, error) {
	emitter := Emitter{
		connection: conn,
	}

	channel, err := emitter.connection.Channel()
	if err != nil {
		return Emitter{}, err
	}
	defer channel.Close()

	err = declareExchange(channel)
	if err != nil {
		return Emitter{}, err
	}

	return emitter, nil
}
