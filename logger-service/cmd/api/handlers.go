package main

import (
	"fmt"
	"net/http"

	"github.com/Nico2220/logger-service/data"
	"github.com/Nico2220/tools"
)

type jsonResponse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	Data any `json:"data,omitempty"`
}


func(app *application) writeLog(w http.ResponseWriter, r *http.Request){
	fmt.Println("Hello from logger New")
	var input struct {
		Name      string `bson:"name" json:"name"`
		Data      string `bson:"data" json:"data"`
	}

	err := tools.ReadJSON(w, r, &input)
	if err != nil {

		payload := jsonResponse{
			Error:   true,
			Message: err.Error(),
		}
		tools.WriteJSON(w, http.StatusBadRequest, payload, nil)
		return
	}

	logEntry := data.LogEntry {
		Name: input.Name,
		Data: input.Data,
	}

	err= app.models.LogEntry.Insert(logEntry)

	if err != nil {
		payload := jsonResponse{
			Error:   true,
			Message: err.Error(),
		}
		tools.WriteJSON(w, http.StatusBadRequest, payload, nil)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: "Logentry inserted with success",
		Data: input.Data,
	}

	tools.WriteJSON(w, http.StatusCreated, payload, nil)

}