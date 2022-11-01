package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/arnab4477/Parkour_API/internal/data"
)

// Handler method on the app instance for the POST /movements endpount
func (app *application) createMovementHandler(w http.ResponseWriter, r *http.Request) {
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
		Prerequisite []string `json:"prerequisite"` 
	}

	// Decode the response body from JSON to a native Go object
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.logError(r, err)
		app.writeError(w, r, http.StatusBadRequest, err.Error())

	}
	
	// Print the input
	fmt.Fprintf(w, "%+v\n", input)
}

// Handler method on the app instance for the GET /movements/:id endpount
func (app *application) showMovementHandler(w http.ResponseWriter, r *http.Request) {

	// Get the id parameter
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		app.logError(r, err)
		app.notFoundResponse(w, r)
		return
	}

	// This is to be used to get the time of UTC instead of local
	utc, _ := time.LoadLocation("UTC")

	// Create a dummy instance of a movement struct
	movement := data.Movement {
		ID: id,
		CreatedAt: time.Now().In(utc),
		Name: "Pull up",
		Description: "Pull up is an essential basic vertical pulling movement",
		Image: "https://i.ytimg.com/vi/HRV5YKKaeVw/maxresdefault.jpg",
		Tutorials: []string{"https://www.youtube.com/watch?v=Y3ntNsIS2Q8", "https://www.gymnasticbodies.com/your-perfect-pull-up/"},
		Skilltype: []string{"basics", "strength", "hypertrophy"},
		Muscles: []string{"lats", "biceps", "forearms"},
		Difficulty: "Beginner",
		Equipments: []string{"pull up bar", "gymnastics rings"},
		Prerequisite: nil,
		Version: 1,
	}

	// Convert the struct to JSON and add that to the response body
	err = app.writeJSON(w, envelope{"movement": movement}, http.StatusOK, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}