package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type ChirpParams struct {
	Body string `json:"body"`
}

type ValidationRes struct {
	Error      string `json:"error,omitempty"`
	ClanedBody string `json:"cleaned_body,omitempty"`
	Body       string `json:"-"`
	Status     int    `json:"-"`
	Valid      bool   `json:"valid"`
}

const (
	charLimit  int    = 140
	tooLong    string = "Chirp is too long"
	unexpected string = "Something went wrong"
)

func handleChirpValidation(w http.ResponseWriter, r *http.Request) {
	// decode req body
	decoder := json.NewDecoder(r.Body)
	chirp := new(ChirpParams)
	err := decoder.Decode(chirp)
	if err != nil {
		log.Printf("ERROR: could not unmarshal reqest, %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}

	// init res
	res := new(ValidationRes)
	res.Body = chirp.Body
	res.checkChirpLength()

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("ERROR: could not marshal response, %s", err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, res.Status, data)
}
