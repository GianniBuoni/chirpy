package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleGETChirpsWithID(
	w http.ResponseWriter, r *http.Request,
) {
	// parse db query
	id, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "uuid.Parse", err)
		return
	}
	data, err := cfg.Queries.GetChirp(r.Context(), id)
	if err != nil {
		respondWithInfoError(w, r.Pattern, http.StatusNotFound)
		return
	}
	// parse response
	res, err := json.Marshal(data)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
