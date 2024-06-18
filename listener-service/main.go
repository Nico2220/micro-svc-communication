package main

import (
	"log"
	"os"

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
}

func connect()(*amqp.Connection, error){

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return nil, err
	} 
	return conn, nil
}