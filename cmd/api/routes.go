package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Route method on the app struct to handle the routing and adding handlers
func (app *application) routes() *httprouter.Router {
	// Declare a new httprouter router instance
	router := httprouter.New()

	// Add the custom error handlers in place of httprouter's default error handler
	// so that it responds with JSON instead of plain text
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// Register the handlers for the endpoints with speific methods
	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/v1/movements", app.createMovementHandler)
	router.HandlerFunc(http.MethodGet, "/v1/movements/:id", app.showMovementHandler)

	return router
}