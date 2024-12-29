package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/db"
	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
	"golang.org/x/crypto/bcrypt"
)

func AddPlacementCoordinator(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	userRole := r.Header.Get("Role")
	if userRole != "ADMIN" {
		http.Error(w, "Unauthorized: Only admins can add placement coordinators", http.StatusUnauthorized)
		return
	}
	
	var coordinator models.PlacementCoordinator
	if err := json.NewDecoder(r.Body).Decode(&coordinator); err!=nil{
		http.Error(w,"Invalid request payload",http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(coordinator.User.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	db := db.InitDB()
	tableName := "admins"
	if utils.CheckTableExists(db, tableName) {
		fmt.Printf("Table '%s' exists.\n", tableName)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName)
		CreateAdminsTable(db)
	}

	adminQuery := `
				INSERT INTO admins(name, password_hash, email, contact, role, created_at)
				VALUES ($1, $2, $3, $4, $5, NOW())
				RETURNING id
			`
	var adminID int
	err = db.QueryRow(adminQuery, coordinator.User.Name, hashedPassword, coordinator.User.Email, coordinator.User.Contact,"placement_coordinator").Scan(&adminID)
	if err != nil {
		log.Printf("Error inserting admin into database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tableName2 := "placement_coordinators"
	if utils.CheckTableExists(db, tableName2) {
		fmt.Printf("Table '%s' exists.\n", tableName2)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName2)
		CreatePlacementCoordinatorsTable(db)
	}
	coordinatorQuery := `
		INSERT INTO placement_coordinators (user_id, usn, branch, batch)
		VALUES ($1, $2, $3, $4)
	`
	_, err = db.Exec(coordinatorQuery, adminID, coordinator.USN, coordinator.Branch, coordinator.Batch)
	if err != nil {
		log.Printf("Error inserting placement coordinator into database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Placement coordinator added successfully"))

}

func CreatePlacementCoordinatorsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS placement_coordinators (
		user_id INT NOT NULL REFERENCES admins(id) ON DELETE CASCADE, -- Foreign key linking to admins table
		usn VARCHAR(15) PRIMARY KEY,          -- USN of the placement coordinator
		branch VARCHAR(50) NOT NULL               -- Branch of the placement coordinator
		batch  INT NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating placement_coordinators table: %v", err)
	} else {
		log.Println("Placement coordinators table ensured to exist.")
	}
}

func GetAllPlacementCoordinators(w http.ResponseWriter, r *http.Request){
	
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}
	query := `
	SELECT pc.usn, a.name, a.email, pc.branch, pc.batch, a.contact, a.created_at 
	FROM placement_coordinators pc
	INNER JOIN admins a ON pc.user_id = a.id;
	`
	rows, err := db.DB.Query(query)
	if err != nil {
		log.Printf("Error fetching placement coordinators: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var coordinators []models.PlacementCoordinator

	for rows.Next() {
		var pc models.PlacementCoordinator
		var admin models.Admin
		if err := rows.Scan(
			&pc.USN,
			&admin.Name,
			&admin.Email,
			&pc.Branch,
			&pc.Batch,
			&admin.Contact,
			&admin.CreatedAt,
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		pc.User = admin
		coordinators = append(coordinators, pc)
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coordinators)
}

func DeletePlacementCoordinator(w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	usn := r.URL.Query().Get("usn")
	if usn == "" {
		http.Error(w, "USN not provided", http.StatusBadRequest)
		return
	}
	var userID int
	query := `SELECT user_id FROM placement_coordinators WHERE usn = $1`
	err := db.DB.QueryRow(query, usn).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Placement coordinator not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error fetching placement coordinator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete the placement coordinator record
	deleteCoordinatorQuery := `DELETE FROM placement_coordinators WHERE usn = $1`
	_, err = db.DB.Exec(deleteCoordinatorQuery, usn)
	if err != nil {
		log.Printf("Error deleting placement coordinator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	deleteAdminQuery := `DELETE FROM admins WHERE id = $1`
		_, err = db.DB.Exec(deleteAdminQuery, userID)
		if err != nil {
			log.Printf("Error deleting admin record: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Placement coordinator deleted successfully"))
}