package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/models"
)

func FilterByBatch(w http.ResponseWriter,r *http.Request){
	if r.Method!=http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
		return
	}

	batch := r.URL.Query().Get("batch")
	if batch == ""{
		http.Error(w,"Batch not provided",http.StatusBadRequest)
		return
	}

	query := `SELECT name, usn, college_email, personal_email,  contact, branch, batch, current_cgpa FROM students WHERE batch = $1`
	
	rows,err := db.DB.Query(query,batch)
	
	if err!=nil{
		log.Printf("Error querying database: %v", err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var students []models.User
	for rows.Next(){
		var student models.User
		if err := rows.Scan(&student.Name, &student.USN, &student.College_Email, &student.Personal_Email, &student.Contact, &student.Branch, &student.Batch, &student.CurrentCGPA); err!=nil{
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}
