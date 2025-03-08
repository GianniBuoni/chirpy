package main

import (
	"log"
	"net/http"
)

const (
	port         string = "8080"
	filePathRoot string = "."
)

func main() {
	// init server
	serveMux := new(http.ServeMux)
	server := new(http.Server)
	server.Handler = serveMux
	server.Addr = ":" + port

	// handle files
	serveMux.Handle("/", http.FileServer(http.Dir(filePathRoot)))

	// run program
	log.Printf("🐹 Serving files from %s on port %s", filePathRoot, port)
	log.Fatal(server.ListenAndServe())
}
