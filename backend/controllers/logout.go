package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"Github.com/Aryan-2511/Placement_NIE/utils"
)

// TokenBlacklist is a simple in-memory store for invalidated tokens.
var TokenBlacklist = make(map[string]time.Time)

// InvalidateToken adds the token to the blacklist with an expiration time.
func InvalidateToken(token string, expiration time.Time) {
	TokenBlacklist[token] = expiration
}

// IsTokenInvalid checks if the token is in the blacklist and invalid.
func IsTokenInvalid(token string) bool {
	expiration, exists := TokenBlacklist[token]
	if !exists {
		return false
	}
	// If the current time is past the token's expiration, remove it from the blacklist
	if time.Now().After(expiration) {
		delete(TokenBlacklist, token)
		return false
	}
	return true
}

// LogoutHandler handles user logout and token invalidation
func LogoutHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization token is required", http.StatusUnauthorized)
		return
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
		return
	}
	tokenString := parts[1]

	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
    	http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
    	return
	}

	// Check and convert the expiration time
	expFloat64, ok := claims["exp"].(float64)
	if !ok {
    	http.Error(w, "Invalid token expiration", http.StatusUnauthorized)
    	return
	}
	expiration := time.Unix(int64(expFloat64), 0)

	// Invalidate the token
	InvalidateToken(tokenString, expiration)

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"message": "Successfully logged out"}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding logout response: %v", err)
		http.Error(w, "Error processing logout", http.StatusInternalServerError)
		return
	}
}
