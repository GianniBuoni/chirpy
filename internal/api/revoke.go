package api

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("%s, %s\n", r.URL, err)
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	params := database.RevokeTokenParams{
		Token:     rToken,
		UpdatedAt: time.Now(),
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}
	cfg.Queries.RevokeToken(r.Context(), params)
	w.WriteHeader(http.StatusNoContent)
}
