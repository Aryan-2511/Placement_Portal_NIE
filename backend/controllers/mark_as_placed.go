package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)
func MarkStudentAsPlaced(w http.ResponseWriter, r *http.Request){

	if r.Method != http.MethodPost{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
		return
	}

	// Validate and parse JWT
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}

	// Check the user's role
	role := claims["role"].(string)
	if role != "ADMIN" && role != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized role", http.StatusForbidden)
		return
	}
	var placement models.PlacedStudent
	if err := json.NewDecoder(r.Body).Decode(&placement);err!=nil{
		http.Error(w,"Invalid input", http.StatusBadRequest)
		return
	}
	if placement.PlacementDate.IsZero(){
		placement.PlacementDate = time.Now()
	}
	db := db.InitDB()
	tableName := "placed_students"
	if utils.CheckTableExists(db, tableName) {
		fmt.Printf("Table '%s' exists.\n", tableName)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName)
		CreatePlacedStudentsTable(db)
	}
	// Check if the student exists in the `students` table
	queryCheck := `SELECT COUNT(*) FROM students WHERE usn = $1`
	var count int
	err = db.QueryRow(queryCheck, placement.USN).Scan(&count)
	if err != nil {
		log.Printf("Error checking student existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if count == 0 {
		http.Error(w, "Student not found in the database", http.StatusBadRequest)
		return
	}
	query := `
	INSERT INTO placed_students(usn,name,email,branch,company,package,placement_date,contact)
	VALUES($1,$2,$3,$4,$5,$6,$7,$8)
	ON CONFLICT(usn) DO NOTHING;
	`
	_,err = db.Exec(query,placement.USN,placement.Name,placement.Email,placement.Branch,placement.Company,placement.Package,placement.PlacementDate,placement.Contact)
	if err != nil {
		log.Printf("Error inserting placed student data: %v",err)
		http.Error(w,"Error marking student as placed",http.StatusInternalServerError)
		return 
	}
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("Student marked as placed successfully"))
}
func CreatePlacedStudentsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS placed_students (
    id SERIAL PRIMARY KEY,
    usn VARCHAR(10) UNIQUE NOT NULL REFERENCES students(usn) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    branch VARCHAR(50) NOT NULL,
    company VARCHAR(100) NOT NULL,
    package NUMERIC(10, 2) NOT NULL,
    placement_date TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	contact VARCHAR(15) NOT NULL
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
		http.Error(w, "Unauthorized: Only admins or placement coordinators can delete opportunities", http.StatusUnauthorized)
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
	query := `
			UPDATE placed_students
			SET name = $1, email = $2 , branch = $3, company = $4, package = $5, placement_date = $6, contact = $7
			WHERE usn = $8
			`
	db := db.InitDB()
	_,err := db.Exec(query,
		placed_student.Name,
		placed_student.Email,
		placed_student.Branch,
		placed_student.Company,
		placed_student.Package,
		placed_student.PlacementDate,
		placed_student.Contact,
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