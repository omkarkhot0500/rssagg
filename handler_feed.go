package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/omkarkhot0500/rssagg/internal/database"
	"github.com/google/uuid"
)

// handlerFeedCreate creates a new feed for a logged-in user
// user is already authenticated and passed into this function
func (cfg *apiConfig) handlerFeedCreate(w http.ResponseWriter, r *http.Request, user database.User) {

	// parameters defines the expected JSON body format
	// example input:
	// { "name": "Tech News", "url": "https://example.com/rss" }
	type parameters struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	}

	// Create JSON decoder to read request body
	decoder := json.NewDecoder(r.Body)

	// Empty struct to store decoded JSON
	params := parameters{}

	// Decode JSON into params struct
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Insert feed into database
	// associate it with the authenticated user
	feed, err := cfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),           // generate unique feed ID
		CreatedAt: time.Now().UTC(),     // timestamp
		UpdatedAt: time.Now().UTC(),     // timestamp
		UserID:    user.ID,              // link feed to user
		Name:      params.Name,          // feed name from request
		Url:       params.URL,           // feed URL from request
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create feed")
		return
	}

	// Convert database feed model into API response format
	// send JSON response back to client
	respondWithJSON(w, http.StatusOK, databaseFeedToFeed(feed))
}
