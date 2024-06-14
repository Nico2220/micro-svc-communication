package main

import (
	"net/http"

	"github.com/Nico2220/tools"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *application) errJSON(w http.ResponseWriter, errObj error, status int) error {
	payload := jsonResponse{
		Error:   true,
		Message: errObj.Error(),
	}

	err := tools.WriteJSON(w, status, payload, nil)
	if err != nil {
		return err
	}

	return err

}
