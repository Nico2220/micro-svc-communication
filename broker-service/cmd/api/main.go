package main

import (
	"fmt"
	"log"
	"net/http"
)
const (
	WebPort = "8080"
)

type Config struct{}
func main(){
	app := Config{}

	log.Printf("starting broker service on port %s\n", WebPort)


	srv := &http.Server{
		Addr: fmt.Sprintf(":%s", WebPort),
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}