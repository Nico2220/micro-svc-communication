package main

import (
	"log"
	"os"
	"time"

	"github.com/Nico2220/listener-service/event"
	amqp "github.com/rabbitmq/amqp091-go"
)


func main (){
	conn, err := connectToRabbitmq()
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

	err = consumer.Listen([]string{"log.INFO", "log.ERROR", "log.WARNING"})
	if err != nil {
		log.Println(err)
	}

}

func connectToRabbitmq()(*amqp.Connection, error){
	count := 0
	for {

		conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
		if err != nil {
			log.Println("rabbitmq not ready yet...")
		} else {
			log.Println("connected to rabbitmq from broker service")
			return conn, nil
		}


		if count >= 10 {
			log.Println(err)
			return nil, nil
		}

		time.Sleep(time.Second * 1)

	}
	

}