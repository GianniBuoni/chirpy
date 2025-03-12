package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(
	userID uuid.UUID,
	tokenSecret string,
	expiresIn time.Duration,
) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})
	jwtString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", err
	}
	return jwtString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// parse token string
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.UUID{}, err
	}

	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return uuid.UUID{}, errors.New("ERROR: could not parse claims type")
	}
	// parse expiration
	if time.Now().UTC().After(claims.ExpiresAt.UTC()) {
		return uuid.UUID{}, errors.New("Token expired")
	}
	// parse id
	id, err := uuid.Parse(claims.Subject)
	return id, nil
}
