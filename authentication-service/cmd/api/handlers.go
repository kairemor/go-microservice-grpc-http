package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/tsawler/toolbox"
)


func (app *Config) Authenticate(w http.ResponseWriter, r *http.Request) {
	var tools toolbox.Tools

	var requestPayload struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	err := tools.readJSON(w, r, &requestPayload)
	if err != nil {
		tools.errorJSON(w, err, http.StatusBadRequest)
		return
	}

	// validate the user against the database
	user, err := app.Models.User.GetByEmail(requestPayload.Email)
	if err != nil {
		tools.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	valid, err := user.PasswordMatches(requestPayload.Password)
	if err != nil || !valid {
		tools.errorJSON(w, errors.New("invalid credentials"), http.StatusBadRequest)
		return
	}

	payload := tools.jsonResponse {
		Error: false,
		Message: fmt.Sprintf("Logged in user %s", user.Email),
		Data: user,
	}

	tools.writeJSON(w, http.StatusAccepted, payload)
}