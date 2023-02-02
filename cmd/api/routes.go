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

	// Register the handlers for the /movements/ endpoints
	router.GET("/v1/healthcheck", app.healthcheckHandler)
	router.GET("/v1/movements", app.allowCORS(app.getMovementsHandler))
	router.GET("/v1/movements/:id", app.allowCORS(app.getOneMovementHandler))

	router.POST("/v1/movements", app.allowCORS(app.authenticate(app.requireActivatedUser(app.createMovementHandler))))
	router.PUT("/v1/movements/:id", app.allowCORS(app.authenticate(app.requireActivatedUser(app.updateMovementHandler))))
	router.DELETE("/v1/movements/:id", app.allowCORS(app.authenticate(app.requireActivatedUser(app.deleteMovementHandler))))

	// Register the handlers for the /users/ endpoints
	router.POST("/v1/users", app.allowCORS(app.registerUserHandler))
	router.POST("/v1/users/activate", app.allowCORS(app.activateUserHandler))
	router.POST("/v1/users/login", app.allowCORS(app.loginHandler))

	return router
}
