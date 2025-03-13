package api

import (
	"encoding/json"
	"net/http"

	"github.com/GianniBuoni/chirpy/internal/auth"
)

type refreshRes struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "GetBearerToken", err)
		return
	}
	user, err := cfg.Queries.GetUserFromRefreshToken(
		r.Context(), rToken,
	)
	if err != nil {
		respondWithInfoError(w, r.Pattern, http.StatusUnauthorized)
		return
	}
	if user.RevokedAt.Valid {
		respondWithInfoError(w, r.Pattern, http.StatusUnauthorized)
		return
	}
	// init res
	data := refreshRes{}
	data.Token, err = auth.MakeJWT(user.ID, cfg.SignSecret, JWTDuration)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "MakeJWT", err)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
