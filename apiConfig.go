package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/skodnik/go-contenttype/contenttype"
)

type apiConfig struct {
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
	cfg.fileserverHits.Store(0)
	log.Print(cfg.hits())
	w.Header().Add("Content-Type", contenttype.TextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cfg.hits()))
}
