package api

import (
	"net/http"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *ApiConfig) AuthMiddleware(
	f func(http.ResponseWriter, *http.Request, uuid.UUID),
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// user authentiction
		token, err := auth.GetBearerToken(r.Header)
		if err != nil {
			respondWithInfoError(w, r.Pattern, http.StatusUnauthorized, "not logged in")
			return
		}
		id, err := auth.ValidateJWT(token, cfg.SignSecret)
		if err != nil {
			respondWithInfoError(w, r.Pattern, http.StatusUnauthorized)
			return
		}
		f(w, r, id)
	}
}
