package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// Method on the app instance to retrieve the id parameter
func (app *application) readIDParam(ps httprouter.Params) (int64, error) {

	// Get the id parameter and make sure it is a positive integer
	id, err := strconv.ParseInt(ps.ByName("id"), 10, 64)
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}

	return id, nil
}

// The envelope type
type envelope map[string]interface{}

// Method pn the app instance to create and write JSON response, set status code and headers
func (app *application) writeJSON(w http.ResponseWriter, data envelope, status int, header http.Header) error {

	// Create JSON from data
	json, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	
	// Append a new line at the end for better view
	json = append(json, '\n')

	// Iterate over the given header parameter and set the headers
	for key, value := range(header) {
		w.Header()[key] = value
	}

	// Set the application/json header
	// Write the appropriate status code and the json sresponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write([]byte(json))
	return nil
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, jsonInput interface{}) error {
	// Set the max siz of the input
	maxBytes := 1_048_576
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	// Initialize the JSON decoder and restricting unknown fields
	// This lines and the subsequest error codes need to be deleted if unknown fields are to be allowed
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	// Decode the response body from JSON to a native Go object
	err := decoder.Decode(jsonInput)
	if err != nil {
		// Declare variables for potential error types
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var InvalidUnmarshalError *json.InvalidUnmarshalError

		// Create a switch-case and return the appropriate error message
		switch {
		case errors.Is(err, io.EOF):
			return fmt.Errorf("body must not be empty")

		case errors.Is(err, io.ErrUnexpectedEOF):
			return errors.New("body contains badly formed JSON")
			
		case errors.As(err, &syntaxError):
			return fmt.Errorf("body contains badly formed JSON (at charater %d)", syntaxError.Offset)
			
		case errors.As(err, &unmarshalTypeError):
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON types for field at %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly formed JSON (at charater %d)", unmarshalTypeError.Offset)
		
		// This case is to be deleted to allow unknown fields
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			unknownKey := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return fmt.Errorf("body contains unknown key %s", unknownKey)

		case err.Error() == "http: request body too large":
			return fmt.Errorf("request body cannot be larger than %d bytes", maxBytes)
			
		case errors.As(err, &InvalidUnmarshalError):
			panic(err)

		default:
			return err
		}

	}

	// Check that the body only contains only one JSON value in the request body
	// This is needed because json.Decode() checks only one value at a time
	// This can also be deleted if multiple values are to be allowed
	// But generally it should not be and only one value per request body is to be preffered
	err = decoder.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain a single JSON value")
	}

	return nil
}
	