package main

import (
	"fmt"
	"net/http"
)

// Function that logs error messgaes
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// Error handler that sends error response back as JSON with custom message and http status
func (app *application) writeError(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	err := app.writeJSON(w, envelope{"error" : message}, status, nil)
	if err != nil {
		app.logError(nil, err)
		w.WriteHeader(500)
		app.logger.Println(err)
	}
}

// Handler that sends an error response in the case of an Internal Server Error
func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {

	// Write and send the appropriate error message
	message := "the server encountered a problem and could not proceed with your request"
	app.writeError(w, r, http.StatusInternalServerError, message)
	app.logger.Println(err)
}
// Handler that sends an error response in the case of the method not being allowed
func (app *application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	// Write and send the appropriate error message
	message := fmt.Sprintf("the %s method is not supported for this resource", r.Method)
	app.writeError(w, r, http.StatusMethodNotAllowed, message)
}
// Handler that sends an error response in the case of a Resource Not Found Error
func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	// Write and send the appropriate error message
	message := "the resource could not be found"
	app.writeError(w, r, http.StatusNotFound, message)
}
// Handler that sends an error response in the case of a Edit Conflict Error
func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {

	// Write and send the appropriate error message
	message := "unable to update the record due an edit conflict, please try again"
	app.writeError(w, r, http.StatusConflict, message)
}
// Handler that sends an error response in the case of a Bad Request
func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.writeError(w, r, http.StatusBadRequest, err.Error())
	app.logger.Println(err)
}

// Handler that sends an error response in the case of an failed validation error
func (app *application) failedValidationError(w http.ResponseWriter, r *http.Request, errors map[string]string){
	app.writeError(w, r, http.StatusUnprocessableEntity, errors)
}

// Handler that sends an error response in the case of an user sending invalid authentication credintials
func (app *application) invalidCredentialsResponse(w http.ResponseWriter, r *http.Request) {
	message := "invalid authentication credentials"
	app.writeError(w, r, http.StatusUnauthorized, message)
	}
