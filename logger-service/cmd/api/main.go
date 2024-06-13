package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Nico2220/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	port     = "8082"
	rpcPort  = "5001"
	mongoURI = "mongodb://localhost:27017"
	gRpcPort = "50001"
)

var client *mongo.Client

type Config struct {
	Models data.Models
}

func main() {
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic()
	}

	log.Println("connted to db...")

	client = mongoClient

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := Config{
		Models: data.New(client),
	}

	log.Println("api is running on port", port)
	
	err = app.serve() 
	if err != nil {
		panic(err)
	}
}


func(app *Config) serve() error{
	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func connectToMongo() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(mongoURI)
	clientOptions.SetAuth(options.Credential{
		Username: "admin",
		Password: "password",
	})
	c, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	return c, nil
}
