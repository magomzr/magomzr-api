package handlers

import (
	"github.com/magomzr/magomzr-api/pkg"
)

func GenerateKey(secretKey string) (string, error) {
	token, err := pkg.GenerateKey(secretKey)
	if err != nil {
		return "", err
	}

	return token, nil
}
