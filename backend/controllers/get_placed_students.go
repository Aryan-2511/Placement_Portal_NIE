package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/models"
)

func GetPlacedStudents(w http.ResponseWriter, r *http.Request){
	query := `SELECT id,usn,name,email,branch,company,package,placement_date FROM placed_students`
	rows,err  := db.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching placed students data")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var placedStudents []models.PlacedStudent
	for rows.Next(){
		var student models.PlacedStudent
		if err := rows.Scan(&student.ID, &student.USN, &student.Name, &student.Email, &student.Branch, &student.Company, &student.Package, &student.PlacementDate); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Error fetching data", http.StatusInternalServerError)
			return
		}
		placedStudents = append(placedStudents, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(placedStudents)
}