package main

import (
	"encoding/json"
	"net/http"
)

type jsonRespnse struct {
	Error bool `json:"error"`
	Message string `json:"message"`
	data any `json:"data"`
}

func(app *Config) broker(w http.ResponseWriter, r *http.Request){
	payload := jsonRespnse {
		Error: false,
		Message: "Hello Broker",
	}

	out, _ := json.MarshalIndent(payload, "", "\t")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(out)
}