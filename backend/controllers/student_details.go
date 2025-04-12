package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

func GetStudentDetailsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
			current_cgpa, backlogs, role, isPlaced, resume_link
		FROM students
		WHERE usn = $1
	`

	var student models.User
	err = db.QueryRow(query, usn).Scan(
		&student.Name, &student.USN, &student.DOB, &student.College_Email,
		&student.Personal_Email, &student.Branch, &student.Batch, &student.Address,
		&student.Contact, &student.Gender, &student.Category, &student.Aadhar,
		&student.PAN, &student.Class_10_Percentage, &student.Class_10_Year,
		&student.Class_10_Board, &student.Class_12_Percentage, &student.Class_12_Year,
		&student.Class_12_Board, &student.Current_CGPA, &student.Backlogs, &student.Role, &student.IsPlaced, &student.Resume_link,
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

func EditStudentDetailsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
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

	// Validate the token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "STUDENT" && claims["role"] != "ADMIN" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Extract USN from the request query parameters
	usn := r.URL.Query().Get("usn")
	if usn == "" {
		http.Error(w, "USN is required", http.StatusBadRequest)
		return
	}

	// Parse the request body to get updated details
	var updatedStudent models.User
	if err := json.NewDecoder(r.Body).Decode(&updatedStudent); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Construct the SQL query dynamically based on provided fields
	query := "UPDATE students SET "
	var args []interface{}
	var updates []string

	if updatedStudent.Name != "" {
		updates = append(updates, "name = ?")
		args = append(args, updatedStudent.Name)
	}
	if updatedStudent.DOB != "" {
		updates = append(updates, "dob = ?")
		args = append(args, updatedStudent.DOB)
	}
	if updatedStudent.College_Email != "" {
		updates = append(updates, "college_email = ?")
		args = append(args, updatedStudent.College_Email)
	}
	if updatedStudent.Personal_Email != "" {
		updates = append(updates, "personal_email = ?")
		args = append(args, updatedStudent.Personal_Email)
	}
	if updatedStudent.Contact != "" {
		updates = append(updates, "contact = ?")
		args = append(args, updatedStudent.Contact)
	}
	if updatedStudent.Branch != "" {
		updates = append(updates, "branch = ?")
		args = append(args, updatedStudent.Branch)
	}
	if updatedStudent.Batch != "" {
		updates = append(updates, "batch = ?")
		args = append(args, updatedStudent.Batch)
	}
	if updatedStudent.Address != "" {
		updates = append(updates, "batch = ?")
		args = append(args, updatedStudent.Batch)
	}
	if updatedStudent.Current_CGPA != 0 {
		updates = append(updates, "current_cgpa = ?")
		args = append(args, updatedStudent.Current_CGPA)
	}
	if updatedStudent.Backlogs >= 0 {
		updates = append(updates, "backlogs = ?")
		args = append(args, updatedStudent.Backlogs)
	}
	if updatedStudent.Resume_link != "" {
		updates = append(updates, "resume_link = ?")
		args = append(args, updatedStudent.Resume_link)
	}

	// Ensure there are updates to apply
	if len(updates) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	query += " " + strings.Join(updates, ", ") + " WHERE usn = ?"
	args = append(args, usn)

	// Execute the update query
	_, err = db.Exec(query, args...)
	if err != nil {
		log.Printf("Error updating student details: %v", err)
		http.Error(w, "Error updating student details", http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Student details updated successfully"})
}
