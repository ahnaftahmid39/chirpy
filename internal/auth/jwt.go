package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	if tokenSecret == "" {
		return "", fmt.Errorf("token secret cannot be empty")
	}
	if expiresIn < 0 {
		return "", fmt.Errorf("expire duration cannot be negative, got expiresIn: %v", expiresIn)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		Subject:   userID.String(),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
	})

	return token.SignedString([]byte(tokenSecret))
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(tokenSecret), nil
		})

	if err != nil {
		return uuid.Nil, err
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	expireTime, err := token.Claims.GetExpirationTime()
	if err != nil {
		return uuid.Nil, err
	}

	if expireTime.Compare(time.Now()) == -1 {
		return uuid.Nil, fmt.Errorf("token expired")
	}

	userIdStr, err := token.Claims.GetSubject()
	if err != nil {
		return uuid.Nil, err
	}

	userId, err := uuid.Parse(userIdStr)
	if err != nil {
		return uuid.Nil, err
	}

	return userId, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	bearer := headers.Get("Authorization")
	if bearer == "" {
		return "", fmt.Errorf("authorization header does not exist")
	}
	split := strings.Fields(bearer)
	if len(split) != 2 || split[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header, valid syntax Bearer <token>")
	}
	return split[1], nil
}
