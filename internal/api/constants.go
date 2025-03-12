package api

import "time"

const (
	unexpected string = "Something went wrong"

	// chrip settings
	tooLong   string = "Chirp is too long"
	charLimit int    = 140

	// token settings
	JWTDuration     time.Duration = 1 * time.Hour
	RefreshDuration time.Duration = 1440 * time.Hour
)

var (
	sanitationWords = []string{"kerfuffle", "sharbert", "fornax"}
)
