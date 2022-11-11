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
func (app *application) getMovementsHandler(w http.ResponseWriter, r *http.Request, _ps httprouter.Params) {
	//Create a struct to hold the output values
	var output struct {
		Name string 
		Description string 
		Image string 
		Tutorials []string 
		Skilltype []string 
		Muscles []string 
		Difficulty string 
		Equipments []string 
		Prerequisites []string
		Sort string
		Page int
		Page_size int  
	}

	// Initiate a new Validator instance
	v := validator.NewValidator()

	// Get the query values from the url
	queries := r.URL.Query()

	// Read the queries and put them into the output struct
	output.Name = app.readStrings(queries, "name", "")  
	output.Description = app.readStrings(queries, "description", "")  
	output.Image = app.readStrings(queries, "image", "")  
	output.Difficulty = app.readStrings(queries, "difficulty", "")  

	output.Tutorials = app.readCsv(queries, "tutorials", []string{})
	output.Skilltype = app.readCsv(queries, "skilltype", []string{})
	output.Muscles = app.readCsv(queries, "muscles", []string{})
	output.Equipments = app.readCsv(queries, "equipments", []string{})
	output.Prerequisites = app.readCsv(queries, "prerequisites", []string{})

	output.Sort = app.readStrings(queries, "sort", "id")
	output.Page = app.readInts(queries, "page", 1, v)
	output.Page_size = app.readInts(queries, "page_size", 20, v)

	// Check if the output data is valid
	if !v.NoErrors() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	// Print the output
	fmt.Fprintf(w, "%+v\n", output)
}


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
	// Else send error response back
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

// Handler method on the app instance for the PUT /movements/:id endpount
func (app *application) updateMovementHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

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
		} else if errors.Is(err, data.ErrEditConflict) {
			app.editConflictResponse(w, r)
		} else {
			app.serverErrorResponse(w, r, err)
			return
		}
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
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// Copy the request body to the appropriate fields of the movement record
	movement.Name = input.Name
	movement.Description = input.Description
	movement.Image = input.Image
	movement.Tutorials = input.Tutorials
	movement.Skilltype = input.Skilltype
	movement.Muscles = input.Muscles
	movement.Difficulty = input.Difficulty
	movement.Equipments = input.Equipments
	movement.Prerequisites = input.Prerequisites

	// Initiate a new Validator instance
	v := validator.NewValidator()


	// If there are no errors then proceed
	// Else send error response back
	if data.ValidateMovement(v, movement); !v.NoErrors() {
		app.failedValidationError(w, r, v.Errors)
		return
	}

	// Call the UpdateMovement method on the movement model passing in the validated movement struct
	// This will update an existing record a record in the Movements table in the database
	err = app.models.Movements.UpdateMovement(movement)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// Send response with the movement data
	err = app.writeJSON(w, envelope{"movement": movement}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// Handler method on the app instance for the DELETE /movements/:id endpount
func (app *application) deleteMovementHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// Get the id parameter
	id, err := app.readIDParam(ps)
	if err != nil || id < 1 {
		app.logError(r, err)
		app.notFoundResponse(w, r)
		return
	}

	//Delete the record from the database
	err = app.models.Movements.DeleteMovement(id)
	if err != nil {
		if errors.Is(err, data.ErrNotFound) {
			app.notFoundResponse(w, r)
			return
		} else {
			app.serverErrorResponse(w, r, err)
			return
		}
	}

	// Send response about the success of deletion
	err = app.writeJSON(w, envelope{"movement": "movement successfully deleted"}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}