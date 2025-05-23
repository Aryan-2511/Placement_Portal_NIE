package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/utils"
)
// GetRecentOpportunities counts opportunities from last 30 days
// Used for student dashboard activity metrics
func GetRecentOpportunities(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

    // Validate the token
    claims, err := utils.ValidateToken(tokenString)
    if err != nil {
        log.Print(err)
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }
    if claims["role"] != "STUDENT" {
        http.Error(w, "Unauthorized access", http.StatusForbidden)
        return
    }

    batch := r.URL.Query().Get("batch")
    if batch == "" {
        http.Error(w, "Batch parameter is required", http.StatusBadRequest)
        return
    }

    query := `
        SELECT COUNT(*)
        FROM opportunities 
        WHERE batch = $1 AND created_at >= NOW() - INTERVAL '30 days'
    `
    rows, err := db.Query(query, batch)
    if err != nil {
        log.Printf("Error fetching recent opportunities: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

	var count int
    err = db.QueryRow(query,batch).Scan(&count)
    if err != nil {
        log.Printf("Error fetching active opportunities count: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    response := map[string]int{"recent_opportunities": count}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetActiveOpportunitiesCount retrieves count of current opportunities
// Filters by batch and active status
func GetActiveOpportunitiesCount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

    // Validate the token
    claims, err := utils.ValidateToken(tokenString)
    if err != nil {
        log.Print(err)
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }
    if claims["role"] != "STUDENT" {
        http.Error(w, "Unauthorized access", http.StatusForbidden)
        return
    }
	batch := r.URL.Query().Get("batch")
    if batch == "" {
        http.Error(w, "Batch parameter is required", http.StatusBadRequest)
        return
    }

    query := `SELECT COUNT(*) FROM opportunities WHERE status = 'ACTIVE' AND batch = $1`
    var count int
    err = db.QueryRow(query,batch).Scan(&count)
    if err != nil {
        log.Printf("Error fetching active opportunities count: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    response := map[string]int{"active_opportunities": count}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}

// GetPlacedStudentsCount tracks placement progress in batch
// Used for student dashboard statistics
func GetPlacedStudentsCount(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

    // Validate the token
    claims, err := utils.ValidateToken(tokenString)
    if err != nil {
        log.Print(err)
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }
    if claims["role"] != "STUDENT" {
        http.Error(w, "Unauthorized access", http.StatusForbidden)
        return
    }

    batch := r.URL.Query().Get("batch")
    if batch == "" {
        http.Error(w, "Batch parameter is required", http.StatusBadRequest)
        return
    }

    query := `SELECT COUNT(*) FROM placed_students WHERE batch = $1`
    var count int
    err = db.QueryRow(query, batch).Scan(&count)
    if err != nil {
        log.Printf("Error fetching placed students count for batch %s: %v", batch, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    response := map[string]int{"placed_students": count}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
func GetTotalApplicationsByStudent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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

    // Validate the token
    claims, err := utils.ValidateToken(tokenString)
    if err != nil {
        log.Print(err)
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }
    if claims["role"] != "STUDENT" {
        http.Error(w, "Unauthorized access", http.StatusForbidden)
        return
    }

	studentUSN := r.URL.Query().Get("usn")
    if studentUSN == "" {
        http.Error(w, "USN is required", http.StatusBadRequest)
        return
    }

    // Query the database for the total applications by the student
    query := `SELECT COUNT(*) FROM applications WHERE student_usn = $1`
    var count int
    err = db.QueryRow(query, studentUSN).Scan(&count)
    if err != nil {
        log.Printf("Error fetching total applications for student %s: %v", studentUSN, err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    // Respond with the total applications count
    response := map[string]int{"total_applications": count}
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
