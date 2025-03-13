package api

import (
	"log"
	"net/http"

	"github.com/skodnik/go-contenttype/contenttype"
)

func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithInfoError(w, r.Pattern, http.StatusForbidden)
		return
	}
	// reset hits
	cfg.FileserverHits.Store(0)
	log.Print(cfg.hits())

	// reset users
	err := cfg.Queries.DeleteUsers(r.Context())
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "db.DeleteUsers", err)
		return
	}
	// form res
	w.Header().Add("Content-Type", contenttype.TextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cfg.hits()))
}
