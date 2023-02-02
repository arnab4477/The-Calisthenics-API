package main

import (
	"context"
	"net/http"

	"github.com/arnab4477/Parkour_API/internal/data"
)

// contextKey type for the key
type contextKey string

// Constamt to hold the string "user" as contextKey
const userContextKey = contextKey("user")

// Method to set the user context to the request context
func (app *application) contextSetUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

// Method to retrieve the User struct from the request context
func (app *application) contextGetUser(r *http.Request) *data.User {
	user, ok := r.Context().Value(userContextKey).(*data.User)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
