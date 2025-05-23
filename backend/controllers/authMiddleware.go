package controllers

import (
	"context"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/utils"
)

type contextKey string

// UserContextKey stores user information in request context
const UserContextKey = contextKey("user")

// AuthMiddleware validates JWT tokens and injects user claims into request context
// Ensures all protected routes have valid authentication
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
