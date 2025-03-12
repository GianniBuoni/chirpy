package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) GETchirps(w http.ResponseWriter, r *http.Request) {
	data, err := cfg.Queries.GetChirps(r.Context())
	if err != nil {
		log.Printf("ERROR: issue getting chirps, %s\n", err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: issue marshaling chirps, %s\n", err)
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) GETchirp(w http.ResponseWriter, r *http.Request) {
	// parse db query
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		log.Printf("Could not parse UUID: %s\n", r.PathValue("chirpID"))
	}
	data, err := cfg.Queries.GetChirp(r.Context(), id)
	if err != nil {
		log.Printf("ERROR: could not make db query, %s\n", err)
		respondWithError(w, http.StatusNotFound, "Chirp not found")
		return
	}
	// parse response
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: could not marshal chirp data, %s\n", err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) HandleChirp(w http.ResponseWriter, r *http.Request) {
	// parses logged in
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Not logged in")
		return
	}
	id, err := auth.ValidateJWT(token, cfg.SignSecret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid credentils")
		return
	}
	// chirp
	decoder := json.NewDecoder(r.Body)
	chipParams := database.CreateChirpParams{}
	err = decoder.Decode(&chipParams)
	if err != nil {
		log.Printf("ERROR: could unmarshal chirp request, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}

	// validation
	err = checkChirpLength(chipParams.Body)
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	chipParams.Body = cleanBody(chipParams.Body)

	// response
	chipParams.ID = uuid.New()
	chipParams.UserID = id
	chipParams.CreatedAt = time.Now()
	chipParams.UpdatedAt = time.Now()
	data, err := cfg.Queries.CreateChirp(r.Context(), chipParams)
	if err != nil {
		log.Printf("ERROR: could not create new chirp, %s.\n", err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: could not marshal new chirp, %s.\n", err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, http.StatusCreated, res)
}

// helpers
func checkChirpLength(s string) error {
	if utf8.RuneCountInString(s) > charLimit {
		return errors.New(tooLong)
	}
	return nil
}

func cleanBody(s string) string {
	body := s
	for _, werd := range sanitationWords {
		re := regexp.MustCompile(`(?i)` + werd)
		body = re.ReplaceAllString(body, "****")
	}
	return body
}
