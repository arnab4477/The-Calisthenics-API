package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/arnab4477/Parkour_API/internal/data"
	"github.com/arnab4477/Parkour_API/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Handler method on the app instance for the POST /users endpount
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) { 
	// Struct ot hold the input for the user's data
	var input struct {
		Username string `json:"username"`
		Email string `json:"email"`
		Password string `json:"password"`
	}

	// Parse the request body to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the input data to the user struct
	user := &data.User{
		Username: input.Username,
		Email: input.Email,
		Activated: false,
	}

	// Hash the user's plain text password
	err = user.Password.SetHash(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Validate the user struct
	v := validator.NewValidator()
	if data.ValidateUser(v, user); !v.NoErrors() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	// Insert the user data into the database
	err = app.models.Users.InsertOneUser(user)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrDuplicateEmail):
				v.AddError("email", "must be unique")
				app.failedValidationError(w, r, v.Errors)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}

	// Generate an activation token for the user
	token, err := app.models.Tokens.NewToken(user.ID, 2*24*time.Hour, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	responseData := map[string]interface{}{
		"activationToken": token.PlainText,
		"user": user,
	}
	
	// If there has been no error, send the user as JSON response to the client
	// along with the appropriate status code
	err = app.writeJSON(w, envelope{"response": responseData}, http.StatusCreated, nil)
	if err !=  nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler method on the app instance for the POST /users/activate endpount
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	// Struct to hold the input
	var input struct {
		TokenPlaintext string `json:"token"`
	}
	// Read the request body and add it to the input struct
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Validate the plaintext token
	v := validator.NewValidator()
	if data.ValidateTokenPlainText(v, input.TokenPlaintext); !v.NoErrors() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	// Retrieve the user from the token
	user, err := app.models.Users.GetUserFromToken(data.ScopeActivation, input.TokenPlaintext)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrNotFound):
				v.AddError("token", "invalid or expired activation token")
				app.failedValidationError(w, r, v.Errors)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Activate and update the user record
	user.Activated = true
	err = app.models.Users.UpdateOneUser(user)
	if err != nil {
		switch {
			case errors.Is(err, data.ErrEditConflict):
				app.editConflictResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
		}
		return
	}
	// Delete all activation tokens for the user
	err = app.models.Tokens.DeleteTokens(user.ID, data.ScopeActivation)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send the updated user details to the client in a JSON response.
	err = app.writeJSON(w, envelope{"user": user}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

