package pkg

import (
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	userSecretEnv = os.Getenv("UserSecretKey")
	audience      = os.Getenv("Audience")
	issuer        = os.Getenv("Issuer")
)

func GenerateKey(userSecretKey string) (string, error) {
	if userSecretEnv == "" {
		log.Printf("The userSecret env variable is empty")
		return "", errors.New("please provide a user secret variable")
	}

	if userSecretKey != userSecretEnv {
		log.Printf("The userSecret is either expired or invalid")
	}

	claims := jwt.MapClaims{
		"name": "Mario Gómez's Blog",
		"exp":  time.Now().Add(time.Minute * 300).Unix(),
		"aud":  audience,
		"iss":  issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("UserSecretKey")))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func ValidateToken(authHeader string) (bool, error) {
	token := strings.TrimPrefix(authHeader, "Bearer ")
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {
		return []byte(userSecretEnv), nil
	})
	if err != nil || !parsedToken.Valid {
		return false, err
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return false, nil
	}

	if claims["aud"] != audience {
		return false, nil
	}

	return true, nil
}
