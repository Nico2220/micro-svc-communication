package main

import (
	"fmt"
	"net/http"
)


func(app *application) logger(w http.ResponseWriter, r *http.Request){
	fmt.Println("Hello from logger")
}