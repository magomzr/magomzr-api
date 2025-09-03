package handlers

import (
	i "github.com/magomzr/magomzr-api/internal"
)

func GenerateKey(secretKey string) (string, error) {
	token, err := i.GenerateKey(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
