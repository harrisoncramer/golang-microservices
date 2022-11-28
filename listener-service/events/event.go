package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(channel *amqp.Channel) error {
	return channel.ExchangeDeclare(
		"logs_topic", // Name of the exchange
		"topic",      // Type
		true,         // Durable?
		false,        // Auto-deleted?
		false,        // internal?
		false,        // no-wait?
		nil,          // arguments
	)
}

func declareRandomQueue(ch *amqp.Channel) (amqp.Queue, error) {
	return ch.QueueDeclare(
		"",    // Name
		false, // durable?
		false, // delete when unused?
		true,  // exclusive? (yes, for current operations)
		false, // no wait?
		nil,   // arguments
	)
}
