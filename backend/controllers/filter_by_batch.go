package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/utils"
)

// Add this struct at the top with other imports
// FilteredStudent represents essential student information for batch-wise filtering
type FilteredStudent struct {
	Name           string  `json:"name"`
	USN            string  `json:"usn"`
	College_Email  string  `json:"college_email"`
	Personal_Email string  `json:"personal_email"`
	Branch         string  `json:"branch"`
	Batch          string  `json:"batch"`
	Contact        string  `json:"contact"`
	Current_CGPA   float64 `json:"current_cgpa"`
}

// FilterByBatch retrieves students from a specific batch
// Returns filtered student details including academic performance
func FilterByBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
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

	// Validate the token
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

	batch := r.URL.Query().Get("batch")
	if batch == "" {
		http.Error(w, "Batch not provided", http.StatusBadRequest)
		return
	}

	query := `SELECT name, usn, college_email, personal_email, branch, batch, contact, current_cgpa 
              FROM students WHERE batch = $1`
	rows, err := db.Query(query, batch)

	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var students []FilteredStudent
	for rows.Next() {
		var student FilteredStudent
		if err := rows.Scan(
			&student.Name,
			&student.USN,
			&student.College_Email,
			&student.Personal_Email,
			&student.Branch,
			&student.Batch,
			&student.Contact,
			&student.Current_CGPA,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
