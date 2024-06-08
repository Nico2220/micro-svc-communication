package main

import (
	"net/http"
)

func (app *Config) broker(w http.ResponseWriter, r *http.Request) {
	payload := jsonRespnse{
		Error:   false,
		Message: "Hello From broker",
	}

	_ = app.writeJSON(w, http.StatusOK, payload, nil)

}
