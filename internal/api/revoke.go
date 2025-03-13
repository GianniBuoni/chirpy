package api

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleRevoke(w http.ResponseWriter, r *http.Request) {
	rToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "GetBearerToken", err)
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
