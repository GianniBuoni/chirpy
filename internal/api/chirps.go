package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"time"
	"unicode/utf8"

	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleGETChirps(w http.ResponseWriter, r *http.Request) {
	var (
		data []database.Chirp
		err  error
	)

	// check for query params
	s := r.URL.Query().Get("author_id")
	if s != "" {
		id, err := uuid.Parse(s)
		if err != nil {
			respondWithUnexpeted(w, r.Pattern, "uuid.Parse", err)
			return
		}
		data, err = cfg.Queries.GetUserChirps(r.Context(), id)
	} else {
		data, err = cfg.Queries.GetChirps(r.Context())
	}
	if err != nil {
		respondWithInfoError(w, r.Pattern, http.StatusNotFound)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
	}
	respondWithJSON(w, http.StatusOK, res)
}

func (cfg *ApiConfig) HandlePOSTChirp(
	w http.ResponseWriter, r *http.Request, id uuid.UUID,
) {
	// parse chirp
	decoder := json.NewDecoder(r.Body)
	chipParams := database.CreateChirpParams{}
	err := decoder.Decode(&chipParams)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "decoder.Decode", err)
		return
	}

	// validation
	ok := checkChirpLength(chipParams.Body)
	if !ok {
		respondWithInfoError(w, r.Pattern, http.StatusBadRequest, tooLong)
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
		respondWithUnexpeted(w, r.Pattern, "db.CreateChirp", err)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, res)
}

// helpers
func checkChirpLength(s string) bool {
	if utf8.RuneCountInString(s) > charLimit {
		return false
	}
	return true
}

func cleanBody(s string) string {
	body := s
	for _, werd := range sanitationWords {
		re := regexp.MustCompile(`(?i)` + werd)
		body = re.ReplaceAllString(body, "****")
	}
	return body
}
