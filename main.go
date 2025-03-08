package main

import (
	"log"
	"net/http"
)

func main() {
	const (
		port         string = "8080"
		filePathRoot string = "."
	)
	// handlers
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthCheck)
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir(filePathRoot))))

	// init server
	server := new(http.Server)
	server.Handler = mux
	server.Addr = ":" + port

	// run program
	log.Printf("🐹 Serving files from %s on port %s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK\n"))
}
