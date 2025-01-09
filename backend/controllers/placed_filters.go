package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"database/sql"
	"strconv"
	"fmt"
	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

func FilterPlacedByBranch(w http.ResponseWriter, r *http.Request,db *sql.DB,secretKey string){
	if r.Method!=http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
		return
	}
	userRole, err := utils.ExtractRoleFromToken(r, secretKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
		return
	}

	// Check if the user is authorized to add admins
	if userRole != "ADMIN" && userRole!="PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins and PCs can add new admins", http.StatusUnauthorized)
		return
	}


	branch := r.URL.Query().Get("branch")
	if branch == ""{
		http.Error(w,"Branch not provided",http.StatusBadRequest)
		return
	}

	query := `SELECT name, usn, email, contact, branch, company, package  FROM placed_students WHERE branch = $1`
	rows,err := db.Query(query,branch)
	
	if err!=nil{
		log.Printf("Error querying database: %v", err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var students []models.PlacedStudent
	for rows.Next(){
		var student models.PlacedStudent
		if err := rows.Scan(&student.Name, &student.USN, &student.Email, &student.Contact, &student.Branch, &student.Company, &student.Package); err!=nil{
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)


}
func FilterPlacedByCompany(w http.ResponseWriter, r *http.Request,db *sql.DB,secretKey string){
	if r.Method!=http.MethodGet{
		http.Error(w,"Invalid request method",http.StatusMethodNotAllowed)
		return
	}
	userRole, err := utils.ExtractRoleFromToken(r, secretKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
		return
	}

	// Check if the user is authorized to add admins
	if userRole != "ADMIN" && userRole!="PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins and PCs can add new admins", http.StatusUnauthorized)
		return
	}

	company := r.URL.Query().Get("company")
	if company == ""{
		http.Error(w,"Branch not provided",http.StatusBadRequest)
		return
	}

	query := `SELECT name, usn, email, contact, branch, company, package  FROM placed_students WHERE company = $1`
	rows,err := db.Query(query,company)
	
	if err!=nil{
		log.Printf("Error querying database: %v", err)
		http.Error(w,"Internal server error",http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var students []models.PlacedStudent
	for rows.Next(){
		var student models.PlacedStudent
		if err := rows.Scan(&student.Name, &student.USN, &student.Email, &student.Contact, &student.Branch, &student.Company, &student.Package); err!=nil{
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		students = append(students, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)

}
func FilterPlacedByCTC(w http.ResponseWriter, r *http.Request,db *sql.DB,secretKey string){
	if r.Method != http.MethodGet{
		http.Error(w, "Invalid request method",http.StatusMethodNotAllowed)
		return
	}
	userRole, err := utils.ExtractRoleFromToken(r, secretKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
		return
	}
	// Check if the user is authorized to add admins
	if userRole != "ADMIN" && userRole!="PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins and PCs can add new admins", http.StatusUnauthorized)
		return
	}

	minCTCParam := r.URL.Query().Get("min_ctc")
	maxCTCParam := r.URL.Query().Get("max_ctc")
	order := r.URL.Query().Get("order")

	if order!= "asc" && order!="desc"{
		order = "asc"
	}

	minCTC := 0
	maxCTC := int(^uint(0)>>1)
	var err1 error

	if minCTCParam != ""{
		minCTC, err1 = strconv.Atoi(minCTCParam)
		if err1 != nil{
			http.Error(w, "Invalid min_ctc value", http.StatusBadRequest)
			return
		}
	}
	if maxCTCParam != ""{
		maxCTC, err1 = strconv.Atoi(maxCTCParam)
		if err1 != nil{
			http.Error(w, "Invalid max_ctc value", http.StatusBadRequest)
			return
		}
	}

	query := `SELECT name, usn, email, contact, branch, company, package
			FROM placed_students WHERE package BETWEEN $1 AND $2
			ORDER BY package` + order
	rows,err1 := db.Query(query,minCTC,maxCTC)
	if err1!=nil{
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return 
	}
	defer rows.Close()

	var students []models.PlacedStudent
	for rows.Next(){
		var student models.PlacedStudent
		if err1:= rows.Scan(&student.Name, &student.USN, &student.Email, &student.Contact, &student.Branch, &student.Company, &student.Package ); err1!=nil{
			log.Printf("Error scanning row: %v", err1)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return 
		}
		students = append(students, student)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)

}
func FilterPlacedHandler(w http.ResponseWriter, r *http.Request,db *sql.DB,secretKey string){
	if r.Method != http.MethodGet{
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	userRole, err := utils.ExtractRoleFromToken(r, secretKey)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unauthorized: %v", err), http.StatusUnauthorized)
		return
	}

	// Check if the user is authorized to add admins
	if userRole != "ADMIN" && userRole!="PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins and PCs can add new admins", http.StatusUnauthorized)
		return
	}


	if branch := r.URL.Query().Get("branch");branch!=""{
		FilterPlacedByBranch(w,r,db,secretKey)
		return
	} 

	if company := r.URL.Query().Get("company");company!=""{
		FilterPlacedByCompany(w,r,db,secretKey)
		return
	}
	 
	if r.URL.Query().Has("min_ctc")||r.URL.Query().Has("max_ctc") {
		FilterPlacedByCTC(w,r,db,secretKey)
		return
	}
	http.Error(w, "Invalid or missing filter parameters", http.StatusBadRequest)
}

