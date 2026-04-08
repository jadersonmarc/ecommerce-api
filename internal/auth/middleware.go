package auth

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserContextKey = contextKey("user")

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return secret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		claims := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), UserContextKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func MiddlewareAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		claims, ok := r.Context().Value(UserContextKey).(jwt.MapClaims)

		if !ok || claims["role"] != "admin" {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
