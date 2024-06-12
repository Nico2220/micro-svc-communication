package main

import (
	"fmt"
	"net/http"
)

func (app *Config) auth(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Hello From Auth-service")
	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errJSON(w, err, http.StatusBadRequest)
		return
	}

	// user, err := app.Models.User.GetByEmail(input.Email)
	// if err != nil {
	// 	app.errJSON(w, errors.New("invalid credantials"), http.StatusBadRequest)
	// 	return 
	// }

	// valid , err := user.PasswordMatches(input.Password)
	// if err!= nil || !valid {
	// 	app.errJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
	// }

	payload := jsonRespnse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", input.Email),
		Data: input,
	}

	app.writeJSON(w, http.StatusOK, payload, nil)

	
}
