package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Nico2220/broker/event"
	"github.com/Nico2220/tools"
)

type RequestPayload struct{
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
	Log LogPayload `json:"log,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"-"`
}

type LogPayload struct {
	Name string `json:"name"`
	Data any `json:"data"` 
}

// { "action":"auth", "auth": {"email":"admin@gmail.com", "password":"admin123"}}

func (app *application) broker(w http.ResponseWriter, r *http.Request) {}

func(app *application) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload) 
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
	}

	
	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w, requestPayload.Auth)
	case "log":
		// app.WriteLogItem(w, requestPayload.Log)
		app.logEventViaRabbitmq(w, requestPayload.Log)
	default:
		app.errJSON(w, errors.New("unknown action"), http.StatusMethodNotAllowed)
	}
}
 
func(app *application) WriteLogItem(w http.ResponseWriter, l LogPayload){
	jsonData, _ := json.Marshal(l)

	request, err := http.NewRequest("POST", fmt.Sprintf("http://logger-service:%d", app.config.port), bytes.NewBuffer(jsonData))
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	client := &http.Client{}
	response, err :=  client.Do(request)
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	defer response.Body.Close()

	var jsonFromLog jsonResponse 

	err = json.NewDecoder(response.Body).Decode(&jsonFromLog)
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromLog.Error {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	payload := jsonResponse {
		Error: jsonFromLog.Error,
		Message: jsonFromLog.Message,
		Data: jsonFromLog.Data,
	}

	tools.WriteJSON(w, response.StatusCode, payload, nil)

}


func(app *application) Authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.Marshal(a)

	request, err := http.NewRequest("POST",  fmt.Sprintf("http://auth-service:%d/auth", app.config.port), bytes.NewBuffer(jsonData))
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return

	}

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return

	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusUnauthorized {
		app.errJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	var jsonFromService jsonResponse

	err =  json.NewDecoder(response.Body).Decode(&jsonFromService)
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	if jsonFromService.Error {
		app.errJSON(w, err, http.StatusUnauthorized)
		return
	}

	payload := jsonResponse {
		Error: false,
		Message: "Authenticated",
		Data: jsonFromService.Data,
	}


	err = tools.WriteJSON(w, http.StatusAccepted, payload, nil)
	if err != nil {
		app.errJSON(w, err, http.StatusInternalServerError)
		return
	}

}

func(app *application) logEventViaRabbitmq(w http.ResponseWriter, l LogPayload){
	err := app.pushToQueue(l.Name, (l.Data).(string))
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	jsonResponse := jsonResponse {
		Error: false,
		Message: "logged via rabbit mq",
	}

	tools.WriteJSON(w, http.StatusAccepted, jsonResponse, nil)
}

func(app *application) pushToQueue(name, msg string) error{
	emitter, err := event.NewEmitter(app.config.rabbitmq)
	if err != nil {
		return err
	}

	payload := LogPayload {
		Name: name,
		Data: msg,
	}


	j, _ := json.Marshal(&payload)

	err = emitter.Push(string(j), "log.INFO")
	if err != nil {
		return err
	}

	return nil
}
