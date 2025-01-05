package controllers

import(
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"Github.com/Aryan-2511/Placement_NIE/models"
)

func FilterByBranch(w http.ResponseWriter,r *http.Request,db *sql.DB){
	if r.Method!=http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
		return
	}

	branch := r.URL.Query().Get("branch")
	if branch == ""{
		http.Error(w,"Branch not provided",http.StatusBadRequest)
		return
	}

	query := `SELECT name, usn, college_email, personal_email, contact, branch, batch, current_cgpa FROM students WHERE branch = $1`
	rows,err := db.Query(query,branch)
	
	if err!=nil{
		log.Printf("Error querying database: %v", err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var students []models.User
	for rows.Next(){
		var student models.User
		if err := rows.Scan(&student.Name, &student.USN, &student.Name, &student.College_Email,&student.Personal_Email, &student.Contact, &student.Branch, &student.Batch, &student.CurrentCGPA); err!=nil{
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}