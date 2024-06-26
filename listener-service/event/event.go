package event

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func declareExchange(ch *amqp.Channel) error {
	return ch.ExchangeDeclare(
		"logs_topic", //name
		"topic", // kind
		true, // durable
		false, // auto-delete
		false, // internal
		false, // noWait
		nil,  // arguments
	)
}

func delcareRandomQueue(ch *amqp.Channel)(amqp.Queue, error) {
	  return ch.QueueDeclare(
		"",   //name
		false, // durable ?
		false, // auto-delete?
		true, //exlusive
		false, //no-waits
		nil, // arguments
	  )
}