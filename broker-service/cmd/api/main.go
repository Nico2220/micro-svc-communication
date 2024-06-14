package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

type config struct{
	port int
	env string
}

type application struct {
	config config
}

func main() {
	var cfg config

	flag.IntVar(&cfg.port, "port", 8080, "API server port")
	flag.StringVar(&cfg.env, "env","dev", "environment")

	app := application {
		config: cfg,
	}

	log.Printf("starting broker service on port %d\n", app.config.port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
