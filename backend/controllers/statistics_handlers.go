package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

// Add authentication checks and CORS handling
// GetTotalStudentsInBatch counts total students in specific batch
// Used for placement statistics and batch management
func GetTotalStudentsInBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Add authentication check
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	batch := r.URL.Query().Get("batch")
	if batch == "" {
		http.Error(w, "Batch parameter is required", http.StatusBadRequest)
		return
	}

	var count int
	query := `SELECT COUNT(*) FROM students WHERE batch = $1`
	err = db.QueryRow(query, batch).Scan(&count)
	if err != nil {
		log.Printf("Error counting students in batch %s: %v", batch, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"total_students": count})
}

// Add similar authentication blocks to other functions
// GetTotalPlacedInBatch counts placed students in specific batch
// Excludes duplicate placements for accurate statistics
func GetTotalPlacedInBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Add same authentication block as above
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	batch := r.URL.Query().Get("batch")
	if batch == "" {
		http.Error(w, "Batch parameter is required", http.StatusBadRequest)
		return
	}

	// Modify the query to be more accurate
	var count int
	query := `
        SELECT COUNT(DISTINCT usn) 
        FROM placed_students 
        WHERE batch = $1`
	err = db.QueryRow(query, batch).Scan(&count)
	if err != nil {
		log.Printf("Error counting placed students in batch %s: %v", batch, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"total_placed": count})
}

func GetEventsToday(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Add same authentication block as above
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	today := time.Now().Format("2006-01-02")
	w.Header().Set("Content-Type", "application/json")

	query := `
        SELECT schedule_id, title, description, start_time, end_time, batch, created_by, created_at 
        FROM schedule 
        WHERE DATE(start_time) = $1 
        ORDER BY start_time`

	rows, err := db.Query(query, today)
	if err != nil {
		log.Printf("Error fetching today's events: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var events []models.Schedule
	for rows.Next() {
		var event models.Schedule
		err := rows.Scan(
			&event.ScheduleID,
			&event.Title,
			&event.Description,
			&event.StartTime,
			&event.EndTime,
			&event.Batch,
			&event.CreatedBy,
			&event.CreatedAt,
		)
		if err != nil {
			log.Printf("Error scanning event: %v", err)
			continue
		}
		events = append(events, event)
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"date":   today,
		"events": events,
	})
}

func GetTotalOpportunitiesForBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Add same authentication block as above
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	batch := r.URL.Query().Get("batch")
	if batch == "" {
		http.Error(w, "Batch parameter is required", http.StatusBadRequest)
		return
	}

	// Modify query to count only active opportunities
	query := `
        SELECT COUNT(*) 
        FROM opportunities 
        WHERE batch = $1 
        AND status = 'ACTIVE'`

	// Add content type header
	w.Header().Set("Content-Type", "application/json")

	var count int
	err = db.QueryRow(query, batch).Scan(&count)
	if err != nil {
		log.Printf("Error counting opportunities for batch %s: %v", batch, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]int{"total_opportunities": count})
}
