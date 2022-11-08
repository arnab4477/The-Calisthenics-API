package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Handler that responds with application status, operating enviroment
// and version in json format. Created on the app instance
func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {

	// Create an envelop instance containing the system information
	env := map[string]string{
				"status": "available",
				"enviroment": app.config.env,
				"version": version,
				}
	
	// Create and send JSON frpm data as the response
	err := app.writeJSON(w, envelope{"system_info": env}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}