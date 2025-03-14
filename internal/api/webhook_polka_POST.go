package api

import (
	"encoding/json"
	"net/http"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

type webhookParams struct {
	Event string `json:"event"`
	Data  struct {
		UserId uuid.UUID `json:"user_id"`
	} `json:"data"`
}

func (cfg *ApiConfig) HandlePolkaWebhook(
	w http.ResponseWriter, r *http.Request,
) {
	paramKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithInfoError(w, r.Pattern, http.StatusUnauthorized, err.Error())
		return
	}
	if paramKey != cfg.PolkaKey {
		respondWithInfoError(w, r.Pattern, http.StatusUnauthorized)
		return
	}
	params := webhookParams{}
	decoder := json.NewDecoder(r.Body)
	err = decoder.Decode(&params)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "decoder.Decode", err)
		return
	}
	if params.Event != "user.upgraded" {
		logInfo(r.Pattern, "Non user upgrade reqest sent to webhook")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	updateParams := database.UpdateRedUserParams{
		ID:          params.Data.UserId,
		IsChirpyRed: true,
	}
	err = cfg.Queries.UpdateRedUser(r.Context(), updateParams)
	if err != nil {
		respondWithInfoError(w, r.Pattern, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
