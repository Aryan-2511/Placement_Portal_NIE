package controllers

import (
	"database/sql"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/db"
)
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "Missing token", http.StatusBadRequest)
		return
	}

	query := `
		UPDATE students
		SET is_verified = TRUE, verification_token = NULL
		WHERE verification_token = $1
		RETURNING college_email
	`
	var email string
	err := db.DB.QueryRow(query, token).Scan(&email)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Email verified successfully. You can now log in."))
}
