package auth

import (
	"errors"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	s := string(b)
	return s, nil
}

func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return err
	}
	return nil
}

func GetBearerToken(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("No Authorization header set")
	}
	bearer := strings.Replace(header, "Bearer", "", 1)
	bearer = strings.Trim(bearer, " ")
	return bearer, nil
}

func GetAPIKey(headers http.Header) (string, error) {
	header := headers.Get("Authorization")
	if header == "" {
		return "", errors.New("No Authorization header set")
	}
	key := strings.Replace(header, "ApiKey", "", 1)
	key = strings.Trim(key, " ")
	return key, nil
}
