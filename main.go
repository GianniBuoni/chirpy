package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/skodnik/go-contenttype/contenttype"
)

func main() {
	const (
		port         string = "8080"
		filePathRoot string = "."
	)

	// load env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	// psql connection
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}
	queries := database.New(conn)

	// api config
	api := new(apiConfig)
	api.queries = queries
	api.platform = platform

	// handlers
	mux := http.NewServeMux()
	mux.Handle("/app/", api.middlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))),
	))
	mux.HandleFunc("GET /admin/metrics", api.handleMetrics)
	mux.HandleFunc("POST /admin/reset", api.handleReset)
	mux.HandleFunc("GET /api/healthz", healthCheck)
	mux.HandleFunc("POST /api/validate_chirp", handleChirpValidation)
	mux.HandleFunc("POST /api/users", api.handeUsers)

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
