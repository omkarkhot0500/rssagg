package auth

import (
	"errors"
	"net/http"
	"strings"
)

// Custom error that we return if Authorization header is missing
// this helps us know exactly what went wrong
var ErrNoAuthHeaderIncluded = errors.New("no authorization header included")

// GetAPIKey reads the API key from HTTP headers
// it expects header format: Authorization: ApiKey <your_key>
func GetAPIKey(headers http.Header) (string, error) {

	// Get the value of Authorization header from request
	authHeader := headers.Get("Authorization")

	// If header is empty → no API key provided
	if authHeader == "" {
		return "", ErrNoAuthHeaderIncluded
	}

	// Split header by space
	// example: "ApiKey 12345"
	// becomes: ["ApiKey", "12345"]
	splitAuth := strings.Split(authHeader, " ")

	// Check if format is correct
	// must have 2 parts and start with "ApiKey"
	if len(splitAuth) < 2 || splitAuth[0] != "ApiKey" {
		return "", errors.New("malformed authorization header")
	}

	// Return the actual API key (second part)
	return splitAuth[1], nil
}
