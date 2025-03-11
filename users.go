package main

import "net/http"

type usersParams struct {
	Email string `json:"email"`
}

type usersRespose struct {
	Error string `json:"error,omitempty"`
}

func handeUsers(w http.ResponseWriter, r http.Request) {
}
