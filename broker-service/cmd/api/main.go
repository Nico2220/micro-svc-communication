package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type config struct{
	port int
	env string
	rabbitmq *amqp.Connection
}

type application struct {
	config config
	
}

func main() {
	var cfg config
	
	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env","dev", "environment")

	flag.Parse()


	conn, err := connectToRabbitmq()
	if err != nil {
		log.Println("connot connect to rabbitmq from the broker")
	}

	cfg.rabbitmq = conn
	app := application {
		config: cfg,
	}

	log.Printf("starting broker service on port %d\n", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
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