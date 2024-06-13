package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type jsonResponse struct {
	Error   bool   `json:"error"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (app *Config) readJSON(w http.ResponseWriter, r *http.Request, data any) error {
	maxBytes := 1048576

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	dec := json.NewDecoder(r.Body)
	err := dec.Decode(&data)
	if err != nil {
		return err
	}

	err = dec.Decode(&struct{}{})

	if err != io.EOF {
		return errors.New("body must have a single json object")
	}

	return nil

}

func (app *Config) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	out, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	out = append(out, '\n')
	if len(headers) > 0 {
		for key, value := range headers {
			w.Header()[key] = value
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_, err = w.Write(out)
	if err != nil {
		return err
	}

	return nil
}

func (app *Config) errJSON(w http.ResponseWriter, errObj error, status int) error {
	payload := jsonResponse{
		Error:   true,
		Message: errObj.Error(),
	}

	err := app.writeJSON(w, status, payload, nil)
	if err != nil {
		return err
	}

	return err

}
