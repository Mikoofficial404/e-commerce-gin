package jwt

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/argon2id"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type MyCustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, err
}

func MakeJWT(userID uuid.UUID, role string, tokenSecret string, expiresIn time.Duration) (string, error) {
	mySigningKey := []byte(tokenSecret)
	claims := MyCustomClaims{
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "ecommerce-access",
			Subject:   userID.String(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(mySigningKey)
	return ss, err
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, string, error) {
	claims := &MyCustomClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})
	if err != nil {
		return uuid.UUID{}, "", fmt.Errorf("failed to parse token: %w", err)
	}
	if !token.Valid {
		return uuid.UUID{}, "", fmt.Errorf("invalid token")
	}
	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		logrus.Infof("Error parsing user ID from claims: %v", err)
		return uuid.Nil, "", err
	}
	return userID, claims.Role, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}
	splitToken := strings.Split(authHeader, " ")
	if len(splitToken) != 2 {
		return "", fmt.Errorf("Error ")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

func GetApiKey(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("no authorization header")
	}
	splitToken := strings.Split(authHeader, " ")

	if len(splitToken) != 2 {
		return "", fmt.Errorf("Error ")
	}
	if splitToken[0] != "ApiKey" {
		return "", fmt.Errorf("invalid authorization header")
	}
	reqToken := strings.TrimSpace(splitToken[1])
	return reqToken, nil
}

func MakeRefreshToken() (string, error) {
	s := make([]byte, 32)
	_, err := rand.Read(s)
	if err != nil {
		return "", err
	}
	encode := hex.EncodeToString(s)
	return encode, nil
}
