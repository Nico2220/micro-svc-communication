package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nico2220/tools"
)

type RequestPayload struct{
	Action string `json:"action"`
	Auth AuthPayload `json:"auth,omitempty"`
}

type AuthPayload struct {
	Email string `json:"email"`
	Password string `json:"-"`
}

// { "action":"auth", "auth": {"email":"admin@gmail.com", "password":"admin123"}}

func (app *application) broker(w http.ResponseWriter, r *http.Request) {
	

}

func(app *application) HandleSubmission(w http.ResponseWriter, r *http.Request) {
	var requestPayload RequestPayload

	err := tools.ReadJSON(w, r, &requestPayload) 
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
	}

	
	switch requestPayload.Action {
	case "auth":
		app.Authenticate(w, requestPayload.Auth)
	default:
		app.errJSON(w, errors.New("unknown action"), http.StatusMethodNotAllowed)
	}
}


func(app *application) Authenticate(w http.ResponseWriter, a AuthPayload) {
	jsonData, _ := json.Marshal(a)

	request, err := http.NewRequest("POST", "http://auth-service:8081/auth", bytes.NewBuffer(jsonData))
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
