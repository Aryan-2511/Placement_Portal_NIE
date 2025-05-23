package controllers

import (
	"database/sql"
	"log"
	"net/http"

)
// VerifyEmailHandler processes email verification tokens
// Activates student account and clears verification token
func VerifyEmailHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	defer func() {
        if r := recover(); r != nil {
            log.Printf("Recovered from panic: %v", r)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
    }()

    token := r.URL.Query().Get("token")
    if token == "" {
        log.Println("Missing or empty token")
        http.Error(w, "Missing token", http.StatusBadRequest)
        return
    }
    log.Printf("Received token: %s", token)

    if db == nil {
        log.Fatal("Database connection is not initialized")
    }

    query := `
        UPDATE students
        SET is_verified = TRUE, verification_token = NULL
        WHERE verification_token = $1
        RETURNING college_email
    `
    var email string
    err := db.QueryRow(query, token).Scan(&email)
    if err == sql.ErrNoRows {
        log.Printf("No matching token found for token: %s", token)
        http.Error(w, "Invalid or expired token", http.StatusBadRequest)
        return
    } else if err != nil {
        log.Printf("Database error: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    if email == "" {
        log.Println("Email not returned by the query")
        http.Error(w, "Email not found", http.StatusInternalServerError)
        return
    }
    log.Printf("Email verified for: %s", email)
    w.Write([]byte("Email verified successfully. You can now log in."))
}
