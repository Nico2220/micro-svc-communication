package event

import (
	"context"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Emitter struct {
	connection *amqp.Connection
}

func NewEmitter(conn *amqp.Connection) (*Emitter, error) {
	emitter := Emitter {
		connection: conn,
	}
	err := emitter.Setup()
	if err != nil {
		return &Emitter{}, err
	}

	return &emitter, nil
}

func(e *Emitter) Setup() error {
	channel, err := e.connection.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()
	return  declareExchange(channel)
	
}

func(e *Emitter) Push(event, severity string) error {
	channel, err := e.connection.Channel()
	if err!=nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second * 5)
	defer cancel()

	err = channel.PublishWithContext(ctx, "logs_topic", severity, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body: []byte(event),
	})

	if err != nil {
		return err
	}


	return nil 
}