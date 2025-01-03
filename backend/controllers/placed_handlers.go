package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/models"
)
func GeneratePlacementID(batch string, serial int) string {
	// Format serial number as a 4-digit string
	serialStr := fmt.Sprintf("%04d", serial)

	// Extract the last two digits of the batch (e.g., "2025" -> "25")
	batchCode := batch[len(batch)-2:]

	// Return formatted Placement-ID
	return fmt.Sprintf("PL%s%s", batchCode, serialStr)
}
func AddPlacedStudent(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// authHeader := r.Header.Get("Authorization")
	// if authHeader == "" {
	// 	http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
	// 	return
	// }
	userRole := r.Header.Get("Role")
	if userRole != "ADMIN" && userRole != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins or placement coordinators can add opportunities", http.StatusUnauthorized)
		return
	}

	var payload struct {
		USN           string `json:"usn"`
		OpportunityID string `json:"opportunity_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if payload.USN == "" || payload.OpportunityID == "" {
		http.Error(w, "USN and OpportunityID are required", http.StatusBadRequest)
		return
	}

	db := db.InitDB()
	defer db.Close()

	// Check if the student exists in the `students` table
	var student struct {
		Name   string
		Email  string
		Branch string
		Batch string 
		Contact string 
	}
	queryStudent := `SELECT name, college_email, branch,batch, contact FROM students WHERE usn = $1`
	err := db.QueryRow(queryStudent, payload.USN).Scan(&student.Name, &student.Email, &student.Branch,&student.Batch, &student.Contact)
	if err != nil {
		log.Printf("Error fetching student details: %v", err)
		http.Error(w, "Student not found in the database", http.StatusBadRequest)
		return
	}

	// Check if the opportunity exists in the `opportunities` table
	var opportunity struct {
		Company string
		Package float64
		Opportunity_type string
	}
	queryOpportunity := `SELECT company, ctc, opportunity_type FROM opportunities WHERE id = $1`
	err = db.QueryRow(queryOpportunity, payload.OpportunityID).Scan(&opportunity.Company, &opportunity.Package, &opportunity.Opportunity_type)
	if err != nil {
		log.Printf("Error fetching opportunity details: %v", err)
		http.Error(w, "Opportunity not found in the database", http.StatusBadRequest)
		return
	}

	// Extract batch from OpportunityID (assume format is "OP2023XXXX")
	if len(payload.OpportunityID) < 6 {
		http.Error(w, "Invalid OpportunityID format", http.StatusBadRequest)
		return
	}
	batch := fmt.Sprintf("20%s", payload.OpportunityID[2:4])

	// Generate Placement ID using the batch
	var serial int
	querySerial := `SELECT COUNT(*) + 1 FROM placed_students WHERE batch = $1`
	err = db.QueryRow(querySerial, batch).Scan(&serial)
	if err != nil {
		log.Printf("Error fetching serial for Placement-ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	placementID := GeneratePlacementID(batch, serial)

	// Insert into placed_students table
	queryInsert := `
	INSERT INTO placed_students(id, usn, opportunity_id, name, email, branch, batch, company, package, placement_date, contact, placement_type)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, CURRENT_TIMESTAMP, $10, $11)
	`
	_, err = db.Exec(queryInsert, placementID, payload.USN, payload.OpportunityID, student.Name, student.Email, student.Branch,student.Batch, opportunity.Company, opportunity.Package, student.Contact, opportunity.Opportunity_type)
	if err != nil {
		log.Printf("Error inserting placed student data: %v", err)
		http.Error(w, "Error marking student as placed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Student marked as placed successfully"))
}
func CreatePlacedStudentsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS placed_students (
    id VARCHAR(15) PRIMARY KEY,
    usn VARCHAR(10) NOT NULL REFERENCES students(usn),
	opportunity_id NOT NULL REFERENCES opportunities(id),
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    branch VARCHAR(50) NOT NULL,
	batch  VARCHAR(10) NOT NULL,
    company VARCHAR(100) NOT NULL,
    package NUMERIC(10, 2) NOT NULL,
    placement_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	contact VARCHAR(15) NOT NULL,
	placement_type VARCHAR(50) DEFAULT 'PLACEMENT'
	);

	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	} else {
		log.Println("Table `placed_students` created or already exists.")
	}
}

func DeletePlacedStudent(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userRole := r.Header.Get("Role")
	if userRole != "ADMIN" && userRole != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins or placement coordinators can delete placed students", http.StatusUnauthorized)
		return
	}
	usn := r.URL.Query().Get("usn")
	if usn == ""{
		http.Error(w,"USN is required", http.StatusBadRequest)
		return
	}
	query := `DELETE FROM placed_students WHERE usn = $1`
	db := db.InitDB()
	result, err := db.Exec(query,usn)
	if err!=nil{
		log.Printf("Error deleting placed student: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "No placed student found with the given USN", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Placed student deleted successfully"))
}

func EditPlacedStudent(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userRole := r.Header.Get("Role")
	if userRole != "ADMIN" && userRole != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins or placement coordinators can edit opportunities", http.StatusUnauthorized)
		return
	}

	var placed_student models.PlacedStudent
	if err := json.NewDecoder(r.Body).Decode(&placed_student); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if placed_student.USN == "" {
		http.Error(w, "USN is required", http.StatusBadRequest)
		return
	}
	db := db.InitDB()
	var exists bool
	checkQuery := `SELECT EXISTS(SELECT 1 FROM placed_students WHERE usn = $1)`
	err := db.QueryRow(checkQuery, placed_student.USN).Scan(&exists)
	if err != nil {
		log.Printf("Error checking if USN exists: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		http.Error(w, "Student not found in the placed_students table", http.StatusNotFound)
		return
	}
	query := `
			UPDATE placed_students
			SET name = $1, email = $2 , branch = $3, batch = $4, company = $5, package = $6, placement_date = $7, contact = $8, placement_type = $9
			WHERE usn = $10
			`
	_,err = db.Exec(query,
		placed_student.Name,
		placed_student.Email,
		placed_student.Branch,
		placed_student.Batch,
		placed_student.Company,
		placed_student.Package,
		placed_student.PlacementDate,
		placed_student.Contact,
		placed_student.PlacementType,
		placed_student.USN,
	)
	if err != nil {
		log.Printf("Error updating placed student: %v", err)
		http.Error(w, "Error updating placed student details", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Placed student details updated successfully"))	
}