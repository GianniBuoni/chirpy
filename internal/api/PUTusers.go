package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) HandlePUTUsers(
	w http.ResponseWriter, r *http.Request, id uuid.UUID,
) {
	// parse params
	decoder := json.NewDecoder(r.Body)
	params := database.UpdateUserParams{
		ID:        id,
		UpdatedAt: time.Now(),
	}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "decoder.Decode", err)
		return
	}
	if params.Email == "" {
		respondWithInfoError(
			w, r.Pattern, http.StatusBadRequest,
			"body with email expected",
		)
		return
	}
	if params.HashedPassword == "" {
		respondWithInfoError(
			w, r.Pattern, http.StatusBadRequest,
			"body with password expected",
		)
		return
	}
	hashedPass, err := auth.HashPassword(params.HashedPassword)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "auth.HashedPassword", err)
		return
	}
	params.HashedPassword = hashedPass
	user, err := cfg.Queries.UpdateUser(r.Context(), params)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "db.UpdateUser", err)
		return
	}

	// int res
	res, err := json.Marshal(user)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
