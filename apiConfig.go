package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/skodnik/go-contenttype/contenttype"
)

type apiConfig struct {
	platform       string
	queries        *database.Queries
	fileserverHits atomic.Int32
}

// getters
func (cfg *apiConfig) hits() string {
	return fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load())
}

// middleware
func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Print(cfg.hits())
		next.ServeHTTP(w, r)
	})
}

// handlers
func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", contenttype.TextHTML)
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load(),
	)
	w.Write([]byte(s))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Endpoint forbidden")
		return
	}
	// reset hits
	cfg.fileserverHits.Store(0)
	log.Print(cfg.hits())

	// reset users
	err := cfg.queries.DeleteUsers(r.Context())
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
