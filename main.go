package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/GianniBuoni/chirpy/internal/api"
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
	// init api
	api := new(api.ApiConfig)

	// load env
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")

	// psql connection
	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Print(err)
		os.Exit(1)
	}

	// config api
	api.Platform = os.Getenv("PLATFORM")
	api.SignSecret = os.Getenv("SIGN_SECRET")
	api.Queries = database.New(conn)

	// handlers
	mux := http.NewServeMux()
	mux.Handle("/app/", api.MiddlewareMetricsInc(
		http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))),
	))
	mux.HandleFunc("GET /admin/metrics", api.HandleMetrics)
	mux.HandleFunc("POST /admin/reset", api.HandleReset)
	mux.HandleFunc("GET /api/healthz", healthCheck)
	mux.HandleFunc("POST /api/users", api.HandeUsers)
	mux.HandleFunc("POST /api/login", api.HandleLogin)
	mux.HandleFunc("POST /api/chirps", api.HandleChirp)
	mux.HandleFunc("GET /api/chirps", api.GETchirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", api.GETchirp)

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
