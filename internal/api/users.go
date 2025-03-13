package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
	"github.com/google/uuid"
)

func (a *ApiConfig) HandeUsers(w http.ResponseWriter, r *http.Request) {
	// decode req
	decoder := json.NewDecoder(r.Body)
	user := database.CreateUserParams{}
	err := decoder.Decode(&user)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "decoder.Decode", err)
		return
	}
	// handle empty fields
	if user.Email == "" {
		respondWithInfoError(w, r.Pattern, http.StatusBadRequest, "expecting email")
		return
	}
	if user.HashedPassword == "" {
		respondWithInfoError(w, r.Pattern, http.StatusBadRequest, "no password set")
		return
	}
	// init response
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.HashedPassword, err = auth.HashPassword(user.HashedPassword)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "HashPassword", err)
		return
	}
	newUser, err := a.Queries.CreateUser(r.Context(), user)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "db.CreateUser", err)
		return
	}
	data, err := json.Marshal(newUser)
	if err != nil {
		respondWithUnexpeted(w, r.Pattern, "json.Marshal", err)
		return
	}
	respondWithJSON(w, http.StatusCreated, data)
}
