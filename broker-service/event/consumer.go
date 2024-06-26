package event

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn *amqp.Connection
	queueName string 
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer :=  Consumer{
		conn: conn,
	}

	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}

	return  consumer, nil
		
}

func(consumer *Consumer) setup() error{
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

type Payload struct {
	Name string `json:"name"`
	Data string `json:"data"`
}

func(consumer *Consumer) Listen(topics []string) error{
	ch, err := consumer.conn.Channel()
	if err != nil {
		return err
	}

	defer ch.Close()

	q, err := delcareRandomQueue(ch)
	if err != nil {
		return err
	}

	for _, s := range topics {
		err = ch.QueueBind(
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

	messages, err := ch.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func(){
		for d := range messages {
			var payload Payload
			_ = json.Unmarshal(d.Body, &payload)
			go handlePayload(payload)

		}
	}()

	fmt.Printf("waiting for echange:logs_topic and queue: %s", q.Name)
	<-forever
	return nil
}

func handlePayload(payload Payload){
	switch payload.Name {
	case "log", "event":
		err := writeLogItem(payload)
		if err != nil {
			log.Println(err)
		}
	case "auth":
		// 
	default:
		fmt.Println("error from default")
	}
}

// func logEvent(entry Payload) error {
// 	fmt.Printf("The entry payload: %+v", entry)
// 	return nil
// }


func writeLogItem(entry Payload) error {
	jsonData, _ := json.Marshal(entry)

	request, err := http.NewRequest("POST", fmt.Sprintf("http://logger-service:%d", 8080), bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	client := &http.Client{}
	response, err :=  client.Do(request)
	if err != nil {
		return err
	}

	defer response.Body.Close()
	fmt.Println("you just make request ot logger-service from writelItem ")

	return nil 

}