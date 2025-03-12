package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/skodnik/go-contenttype/contenttype"
)

// middleware
func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.FileserverHits.Add(1)
		log.Print(cfg.hits())
		next.ServeHTTP(w, r)
	})
}

// handlers
func (cfg *ApiConfig) HandleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", contenttype.TextHTML)
	w.WriteHeader(http.StatusOK)
	s := fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.FileserverHits.Load(),
	)
	w.Write([]byte(s))
}
