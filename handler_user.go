package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/omkarkhot0500/rssagg/internal/database"
)

// handlerUsersCreate handles POST /users request
// it reads user data from request body and saves it in database
func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {

	// parameters struct defines expected JSON input format
	// example input: { "name": "Omkar" }
	type parameters struct {
		Name string `json:"name"`
	}

	// Create JSON decoder to read request body
	decoder := json.NewDecoder(r.Body)

	// Empty struct to store decoded JSON
	params := parameters{}

	// Decode JSON request into params struct
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters")
		return
	}

	// Call database function to insert new user
	// generate UUID and timestamps automatically
	user, err := cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Convert DB user model into API response format
	// then send JSON response back to client
	respondWithJSON(w, http.StatusOK, databaseUserToUser(user))
}
