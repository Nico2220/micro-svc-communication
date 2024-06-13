package main

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


const (
	port = "8080"
	rpcPort = "5001"
	mongoURI = "mongodb://mongo/27027"
	gRpcPort = "50001"
)

var client *mongo.Client

func main(){
	mongoClient, err := connectToMongo()
	if err != nil {
		log.Panic()
	}

	client = mongoClient
}


func connectToMongo() (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20 * time.Second)
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