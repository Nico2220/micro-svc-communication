package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/Nico2220/auth-service/data"
)

const (
	Port = "8080"
)

type Config struct {
	DB     *sql.DB
	Models data.Models
}

func main() {
	app := Config{}

	log.Printf("starting broker service on port %s\n", Port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", Port),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}
