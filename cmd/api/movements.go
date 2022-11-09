package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/arnab4477/Parkour_API/internal/data"
	"github.com/arnab4477/Parkour_API/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Handler method on the app instance for the POST /movements endpount
func (app *application) createMovementHandler(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	// Send an appropriate error response if the medthod is not POST
	if r.Method != http.MethodPost {
		app.methodNotAllowedResponse(w, r)
		return
	}

	//Create a struct to hold the input
	var input struct {
		Name string `json:"name"`
		Description string `json:"description"`
		Image string `json:"image"`
		Tutorials []string `json:"tutorials"` 
		Skilltype []string `json:"skilltype"`
		Muscles []string `json:"muscles"`
		Difficulty string `json:"difficulty"` 
		Equipments []string `json:"equipments"`
		Prerequisites []string `json:"prerequisite"` 
	}

	// Decode the JSON request and send an appropriate response in case of an error
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Create a new movement instance with the input data
	movement := &data.Movement{
		Name: input.Name,
		Description: input.Description,
		Image: input.Image,
		Tutorials: input.Tutorials,
		Skilltype: input.Skilltype,
		Muscles: input.Muscles,
		Difficulty: input.Difficulty,
		Equipments: input.Equipments,
		Prerequisites: input.Prerequisites,
	}

	// Initiate a new Validator instance
	v := validator.NewValidator()


	// If there are no errors then proceed
	// Else senf error response back
	if data.ValidateMovement(v, movement); !v.NoErrors() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	// Call the insert method on the movement model passing in the validated movement struct
	// This will create a new record in the Movements table in the database
	err = app.models.Movements.InsertMovement(movement)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Creating a location header to let the client know where they can find the newly created information
	header := make(http.Header)
	header.Set("Location", fmt.Sprintf("v1/movements/%d", movement.ID))

	// Send a response with the appropriate status code (201), the movement data and the header
	err = app.writeJSON(w, envelope{"movement": movement}, http.StatusCreated, header)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler method on the app instance for the GET /movements/:id endpount
func (app *application) showMovementHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the id parameter
	id, err := app.readIDParam(ps)
	if err != nil || id < 1 {
		app.logError(r, err)
		app.notFoundResponse(w, r)
		return
	}

	// Fetch data for a specific movement
	movement, err := app.models.Movements.GetMovement(id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			app.notFoundResponse(w, r)
			return
		} else {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Send response with the movement data
	err = app.writeJSON(w, envelope{"movement": movement}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}