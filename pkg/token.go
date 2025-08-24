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
	validationToken = os.Getenv("ValidationToken")
	userSecretEnv   = os.Getenv("UserSecretKey")
	audience        = os.Getenv("Audience")
	issuer          = os.Getenv("Issuer")
)

func GenerateKey(tokenFromUser string) (string, error) {
	if validationToken == "" || userSecretEnv == "" {
		log.Printf("Some tokens are missing or invalid")
		return "", errors.New("please provide the user secret variables")
	}

	if tokenFromUser != validationToken {
		log.Printf("The given token is either expired or invalid")
		return "", errors.New("please verify your token")
	}

	claims := jwt.MapClaims{
		"name": "Mario Gómez's Blog",
		"exp":  time.Now().Add(time.Minute * 300).Unix(),
		"aud":  audience,
		"iss":  issuer,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(userSecretEnv))
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
