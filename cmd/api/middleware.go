package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/arnab4477/Parkour_API/internal/data"
	"github.com/arnab4477/Parkour_API/internal/validator"
	"github.com/julienschmidt/httprouter"
)

// Middleware to authenticate an user making the request
func (app *application) authenticate( next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		
		// Add the "Vary: Authorization" header to the response
		w.Header().Add("Vary", "Authorization")

		// Retrieve the value of the Authorization header from the request
		authorizationHeader := r.Header.Get("Authorization")
		// If there is no Authorization header, add the AnonymousUser
		//to the request context
		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next(w, r, ps)
			return
		}

		// Check that the authentication type is "Bearer" and it is properly formed
		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Extract the authentication token and validate it
		token := headerParts[1]
		v := validator.NewValidator()
		if data.ValidateTokenPlainText(v, token); !v.NoErrors() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		// Get the user associated with the tokwn
		user, err := app.models.Users.GetUserFromToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
				case errors.Is(err, data.ErrNotFound):
					app.invalidAuthenticationTokenResponse(w, r)
				default:app.serverErrorResponse(w, r, err)
			}
			return
		}

		// Add the user information to the request context
		r = app.contextSetUser(r, user)
		next(w, r, ps)
	}
}

// Middleware that checks if an user is authorized to make a request to an endpoint
func (app *application) requireActivatedUser(next httprouter.Handle) httprouter.Handle{
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			// Get the user information from the request context
			user := app.contextGetUser(r)

			// Check if the user is not authenicated
			if user.IsAnonymous() {
				app.authenticationRequiredResponse(w, r)
				return
			}

			// Check if the user is not activated
			if !user.Activated {
				app.inactiveAccountResponse(w, r)
				return
			}

			next(w, r, ps)
	}
}

// Middleware for enabling CORS and handle pre-flight request
func (app *application) allowCORS(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		w.Header().Add("Vary", "Origin")
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("Origin")
		if origin != "" {
			// Enable CORS for all origins
			w.Header().Set("Access-Control-Allow-Origin", "*")
			
			// Check if the request is a pre-flight request 
			// If it is, set the necessary headers and send a 200 OK status back
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE")
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")

				w.WriteHeader(http.StatusOK)
				return
			}
		}

		next(w, r, ps)
	}
}