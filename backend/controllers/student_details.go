package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/models"
)

func GetStudentDetailsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Extract USN from the request
	usn := r.URL.Query().Get("usn")
	if usn == "" {
		http.Error(w, "USN is required", http.StatusBadRequest)
		return
	}

	// Query to fetch student details based on USN
	query := `
		SELECT name, usn, dob, college_email, personal_email, branch, batch, address, 
			contact, gender, category, aadhar, pan, class_10_percentage, class_10_year, 
			class_10_board, class_12_percentage, class_12_year, class_12_board, 
			current_cgpa, backlogs, role, isPlaced
		FROM students
		WHERE usn = $1
	`

	var student models.User
	err := db.QueryRow(query, usn).Scan(
		&student.Name, &student.USN, &student.DOB, &student.College_Email,
		&student.Personal_Email, &student.Branch, &student.Batch, &student.Address,
		&student.Contact, &student.Gender, &student.Category, &student.Aadhar,
		&student.PAN, &student.Class_10_Percentage, &student.Class_10_Year,
		&student.Class_10_Board, &student.Class_12_Percentage, &student.Class_12_Year,
		&student.Class_12_Board, &student.Current_CGPA, &student.Backlogs, &student.Role, &student.IsPlaced,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send student details as JSON response
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(student)
	if err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
