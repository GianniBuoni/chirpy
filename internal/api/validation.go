package api

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"unicode/utf8"
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
	res.cleanBody()

	data, err := json.Marshal(res)
	if err != nil {
		log.Printf("ERROR: could not marshal response, %s", err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, res.Status, data)
}

func (v *ValidationRes) checkChirpLength() {
	if utf8.RuneCountInString(v.Body) > charLimit {
		v.Error = tooLong
		v.Status = http.StatusBadRequest
		return
	}
	v.Valid = true
	v.Error = ""
	v.Status = http.StatusOK
}

func (v *ValidationRes) cleanBody() {
	if v.Valid {
		body := v.Body
		for _, werd := range sanitationWords {
			re := regexp.MustCompile(`(?i)` + werd)
			body = re.ReplaceAllString(body, "****")
		}
		v.ClanedBody = body
	}
}
