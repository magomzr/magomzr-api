package main

import (
	"net/http"
	"strings"

	"github.com/magomzr/magomzr-api/pkg"
)

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "authorization header required", http.StatusUnauthorized)
			return
		}

		valid, err := pkg.ValidateToken(authHeader)
		if err != nil || !valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
