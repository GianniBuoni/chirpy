package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	"github.com/skodnik/go-contenttype/contenttype"
)

func main() {
	const (
		port         string = "8080"
		filePathRoot string = "."
	)
	// api config
	api := new(apiConfig)

	// handlers
	mux := http.NewServeMux()
	mux.Handle("/app/", api.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))),
	))
	mux.HandleFunc("GET /admin/metrics", api.handleMetrics)
	mux.HandleFunc("POST /admin/reset", api.handleReset)
	mux.HandleFunc("GET /api/healthz", healthCheck)

	// init server
	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":" + port

	// run program
	log.Printf("🐹 Serving files from %s on port %s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

const ()

type apiConfig struct {
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		log.Print(cfg.hits())
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) hits() string {
	return fmt.Sprintf("Hits: %d\n", cfg.fileserverHits.Load())
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", contenttype.TextHTML)
	w.WriteHeader(http.StatusOK)
	fmt.Sprintf(
		`<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`, cfg.fileserverHits.Load(),
	)
	w.Write([]byte(cfg.hits()))
}

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	log.Print(cfg.hits())
	w.Header().Add("Content-Type", contenttype.TextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cfg.hits()))
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", contenttype.TextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
