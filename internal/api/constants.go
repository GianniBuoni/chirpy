package api

import (
	"net/http"
	"time"
)

const (
	// chrip settings
	tooLong   string = "chirp too long"
	charLimit int    = 140

	// token settings
	JWTDuration     time.Duration = 1 * time.Hour
	RefreshDuration time.Duration = 1440 * time.Hour
)

var (
	sanitationWords = []string{"kerfuffle", "sharbert", "fornax"}

	// error messages
	errorMesages = map[int]string{
		http.StatusInternalServerError: "Something went wrong",
		http.StatusUnauthorized:        "Invalid credentials",
		http.StatusNotFound:            "Not Found",
		http.StatusBadRequest:          "Bad request",
		http.StatusForbidden:           "Endpoint forbidden",
	}
)
