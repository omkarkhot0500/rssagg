package main

import (
	"net/http"

	// auth package contains helper to read API key from request headers
	"github.com/omkarkhot0500/rssagg/internal/auth"

	// database package contains DB models and queries
	"github.com/omkarkhot0500/rssagg/internal/database"
)

// authedHandler is a custom function type.
// Any handler that wants authentication must match this shape.
// It receives:
// w -> response writer
// r -> request
// user -> authenticated user from database
type authedHandler func(http.ResponseWriter, *http.Request, database.User)

// middlewareAuth is a middleware function.
// It wraps a handler and runs authentication BEFORE the handler runs.
func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {

	// This returned function is the actual HTTP handler
	return func(w http.ResponseWriter, r *http.Request) {

		// Step 1: Read API key from request headers
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			// If API key missing or invalid → stop request
			respondWithError(w, http.StatusUnauthorized, "Couldn't find api key")
			return
		}

		// Step 2: Use API key to fetch user from database
		user, err := cfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			// If no user found → stop request
			respondWithError(w, http.StatusNotFound, "Couldn't get user")
			return
		}

		// Step 3: Call the real handler
		// We pass the authenticated user into it
		handler(w, r, user)
	}
}
