package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/arnab4477/Parkour_API/internal/data"
	"github.com/arnab4477/Parkour_API/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Handler method for the /user/login endpoint
func (app *application) loginHandler(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	// Inputstruct to hold the email and password
	var input struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// Read the input JSON and add it to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the email and password
	v := validator.NewValidator()
	data.ValidateEmail(v, input.Email)
	data.ValidatePlainPassword(v, input.Password)
	if !v.NoErrors() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	// Lookup the user record based on the email 
	user, err := app.models.Users.GetOneUserByEmail(input.Email)
	if err != nil {
	switch {
		case errors.Is(err, data.ErrNotFound):
			app.invalidCredentialsResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Check if the provided password matches the user's password
	match, err := user.Password.Matchhash(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	// Create the token with 30 days as expiry
	token, err := app.models.Tokens.NewToken(user.ID, (24 * 30)*time.Hour, data.ScopeAuthentication)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send the token as response
	err = app.writeJSON(w, envelope{"authentication": token}, http.StatusCreated, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}