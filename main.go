package main

import (
	"log"
	"net/http"

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
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidation)

	// init server
	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":" + port

	// run program
	log.Printf("🐹 Serving files from %s on port %s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", contenttype.TextPlain)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
