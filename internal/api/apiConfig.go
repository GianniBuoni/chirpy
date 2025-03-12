package api

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/skodnik/go-contenttype/contenttype"
)

type ApiConfig struct {
	Platform       string
	Queries        *database.Queries
	FileserverHits atomic.Int32
}

func NewAPI(p string, q *database.Queries) *ApiConfig {
	api := new(ApiConfig)
	api.Platform = p
	api.Queries = q
	return api
}

// getters
func (cfg *ApiConfig) hits() string {
	return fmt.Sprintf("Hits: %d\n", cfg.FileserverHits.Load())
}

// helpers
func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Add("Content-Type", contenttype.ApplicationJSON)
	w.WriteHeader(code)
	data := []byte(fmt.Sprintf("{\"error\": \"%s\"}", msg))
	w.Write(data)
}

func respondWithJSON(w http.ResponseWriter, code int, payload []byte) {
	w.Header().Add("Content-Type", contenttype.ApplicationJSON)
	w.WriteHeader(code)
	w.Write(payload)
}
