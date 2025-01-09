package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"Github.com/Aryan-2511/Placement_NIE/models"
)

func GetPlacedStudents(w http.ResponseWriter, r *http.Request,db *sql.DB,secretKey string){
	
	userRole := r.Header.Get("Role")
	if userRole != "ADMIN" && userRole != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins or placement coordinators can edit opportunities", http.StatusUnauthorized)
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