package main

import "net/http"

func main() {
	serveMux := new(http.ServeMux)
	server := new(http.Server)
	server.Handler = serveMux
	server.Addr = ":8080"

	server.ListenAndServe()
}
