package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

func ProtectedHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(UserContextKey).(jwt.MapClaims)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	email := claims["email"].(string)
	role := claims["role"].(string)

	response := map[string]string{
		"email": email,
		"role":  role,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
