package main

import (
	"log"
	"os"

	"github.com/Nico2220/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)


func main (){
	conn, err := connect()
	if  err != nil {
		log.Println(err)
		os.Exit(1)
	}
	defer conn.Close()
	
	log.Println("connected to rabbitmq...")

	consumer, err := event.NewConsumer(conn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.Info", "log.Error", "log.Warning"})
	if err != nil {
		log.Println(err)
	}

}

func connect()(*amqp.Connection, error){

	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		return nil, err
	} 
	return conn, nil
}