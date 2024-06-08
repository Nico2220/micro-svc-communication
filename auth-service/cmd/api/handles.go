package main

import (
	"fmt"
	"net/http"
)

func (app *Config) auth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello from auth")
}
