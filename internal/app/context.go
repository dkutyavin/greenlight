package app

import (
	"context"
	"net/http"

	"greenlight.dekutyavin.net/internal/data"
)

type contextKey string

const authenticatedUserContextKey = contextKey("authenticatedUser")

func (app *Application) contextSetAuthenticatedUser(r *http.Request, user *data.User) *http.Request {
	ctx := context.WithValue(r.Context(), authenticatedUserContextKey, user)
	return r.WithContext(ctx)
}

func (app *Application) contextGetAuthenticatedUser(r *http.Request) (*data.User, bool) {
	user, ok := r.Context().Value(authenticatedUserContextKey).(*data.User)
	return user, ok
}
