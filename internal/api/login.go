package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GianniBuoni/chirpy/internal/auth"
	"github.com/GianniBuoni/chirpy/internal/database"
)

func (cfg *ApiConfig) HandleLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := database.User{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("ERROR: could not unmarshal request, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	// handle empty fields
	if params.Email == "" {
		respondWithError(w, http.StatusBadRequest, "Bad request: experting email\n")
		return
	}
	if params.HashedPassword == "" {
		respondWithError(w, http.StatusBadRequest, "Bad request: no password set\n")
		return
	}
	user, err := cfg.Queries.GetUser(r.Context(), params.Email)
	if err != nil {
		msg := fmt.Sprintf("Account for %s, not found\n", params.Email)
		respondWithError(w, http.StatusNotFound, msg)
		return
	}
	err = auth.CheckPasswordHash(user.HashedPassword, params.HashedPassword)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}
	// init response
	data, err := cfg.Queries.GetUserResponse(r.Context(), user.ID)
	if err != nil {
		log.Printf("ERROR: could not get user login response, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	res, err := json.Marshal(data)
	if err != nil {
		log.Printf("ERROR: could not get user login response, %s\n", err.Error())
		respondWithError(w, http.StatusInternalServerError, unexpected)
		return
	}
	respondWithJSON(w, http.StatusOK, res)
}
