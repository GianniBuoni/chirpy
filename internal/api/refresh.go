package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/GianniBuoni/chirpy/internal/auth"
)

type refreshRes struct {
	Token string `json:"token"`
}

func (cfg *ApiConfig) HandleRefresh(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Println(err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	user, err := cfg.Queries.GetUserFromRefreshToken(
		r.Context(), rToken,
	)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invaild credentials")
		return
	}
	if user.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "Invaild credentials")
		return
	}
	// init res
	data := refreshRes{}
	data.Token, err = auth.MakeJWT(user.ID, cfg.SignSecret, JWTDuration)
	if err != nil {
		log.Printf("ERROR: %s making JWT\n", r.URL)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: %s marshaling token\n", r.URL)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
