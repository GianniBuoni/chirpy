package api

import (
	"net/http"

	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandleDELETEChirp(
	w http.ResponseWriter, r *http.Request, userId uuid.UUID,
) {
	// parse chirp
	chirpId, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "uuid.Parse", err)
		return
	}
	chirp, err := cfg.Queries.GetChirp(r.Context(), chirpId)
	if err != nil {
		respondWithInfoError(w, r.Pattern, http.StatusNotFound)
		return
	}
	if chirp.UserID != userId {
		respondWithInfoError(w, r.Pattern, http.StatusForbidden)
		return
	}

	// delete chirp
	err = cfg.Queries.DeleteChirp(r.Context(), chirpId)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "db.DeleteChirp", err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
