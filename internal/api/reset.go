package api

import (
	"log"
	"net/http"

	"github.com/skodnik/go-contenttype/contenttype"
)

func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.Platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Endpoint forbidden")
		return
	}
	// reset hits
	cfg.FileserverHits.Store(0)
	log.Print(cfg.hits())

	// reset users
	err := cfg.Queries.DeleteUsers(r.Context())
	if err != nil {
		log.Printf("ERROR: could not reset users, %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}

	// form res
	w.Header().Add("Content-Type", contenttype.TextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cfg.hits()))
}
