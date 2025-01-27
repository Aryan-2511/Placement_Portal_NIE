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

func GetPlacedStudents(w http.ResponseWriter, r *http.Request,db *sql.DB){
	
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
	query := `SELECT id,usn,name,email,branch,batch,company,package,placement_date,contact,placement_type FROM placed_students`
	rows,err  := db.Query(query)
	if err != nil {
		log.Printf("Error fetching placed students data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var placedStudents []models.PlacedStudent
	for rows.Next(){
		var student models.PlacedStudent
		if err := rows.Scan(&student.ID, &student.USN, &student.Name, &student.Email, &student.Branch,&student.Batch, &student.Company, &student.Package, &student.PlacementDate, &student.Contact, &student.PlacementType); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}
		placedStudents = append(placedStudents, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(placedStudents)
}

func GetUnplacedStudents(w http.ResponseWriter, r *http.Request,db *sql.DB){
	
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
	query := `SELECT usn,name,college_email,branch,batch,gender,contact FROM students WHERE students.isPlaced = 'NO'`
	rows,err  := db.Query(query)
	if err != nil {
		log.Printf("Error fetching placed students data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var UnplacedStudents []models.User
	for rows.Next(){
		var student models.User
		if err := rows.Scan(&student.USN, &student.Name, &student.College_Email, &student.Branch,&student.Batch,&student.Gender, &student.Contact); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}
		UnplacedStudents = append(UnplacedStudents, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(UnplacedStudents)
}