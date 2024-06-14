package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Nico2220/logger-service/data"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// const (
// 	port     = "8082"
// 	rpcPort  = "5001"
// 	mongoURI = "mongodb://localhost:27017"
// 	gRpcPort = "50001"
// )



type config struct {
	Models data.Models
	port int
	env string
	rpcPort int
	gRPCPort int
	db struct {
		URI string
	}
}

type application struct {
	config config
	models data.Models
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env", "dev", "evironment")
	flag.IntVar(&cfg.rpcPort, "rpcPort", 5001, "rpc port")
	flag.IntVar(&cfg.gRPCPort, "gRPCPort", 50001, "gRPCPort port")
	flag.StringVar(&cfg.db.URI, "MONGO_URI", "mongodb://mongo:27017", "mongo db uri")
	
	mongoClient, err := connectToMongo(cfg)
	if err != nil {
		log.Panic()
	}

	log.Println("connected to mongodb...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	// close connection
	defer func() {
		if err = mongoClient.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	app := application{
		config: cfg,
		models: data.New(mongoClient),
	}

	log.Println("api is running on port", app.config.port)

	err = app.serve() 
	if err != nil {
		panic(err)
	}
}


func(app *application) serve() error{
	srv := &http.Server{
		Addr: fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func connectToMongo(cfg config) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.db.URI)
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
