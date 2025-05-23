package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
	"golang.org/x/crypto/bcrypt"
)

// AddPlacementCoordinator creates new coordinator with dual table entries
// Links coordinator to both admins and students tables
func AddPlacementCoordinator(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if claims["role"] != "ADMIN" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	var coordinator models.PlacementCoordinator
	if err = json.NewDecoder(r.Body).Decode(&coordinator); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(coordinator.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Insert into database
	tableName := "admins"
	exists, err := utils.CheckTableExists(db, tableName)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Printf("Table '%s' does not exist. Creating table...", tableName)
		CreateAdminsTable(db) // Ensure this function is correctly implemented
	}
	// Insert into database
	tableName2 := "placement_coordinators"
	exists, err = utils.CheckTableExists(db, tableName2)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Printf("Table '%s' does not exist. Creating table...", tableName2)
		CreatePlacementCoordinatorsTable(db) // Ensure this function is correctly implemented
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()
	var serial int
	serialQuery := `SELECT COUNT(*) + 1 FROM admins WHERE role = 'PLACEMENT_COORDINATOR'`
	err = tx.QueryRow(serialQuery).Scan(&serial)
	if err != nil {
		log.Printf("Error fetching serial for Admin-ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Generate Admin ID using the GenerateAdminID function
	adminID := GenerateAdminID("PLACEMENT_COORDINATOR", serial)

	// Insert into admins table
	adminQuery := `
		INSERT INTO admins (id, name, password_hash, email, contact, role, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, NOW())
	`
	_, err = tx.Exec(adminQuery, adminID, coordinator.Name, hashedPassword, coordinator.Email, coordinator.Contact, "PLACEMENT_COORDINATOR")
	if err != nil {
		log.Printf("Error inserting admin: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Insert into placement_coordinators table
	coordinatorQuery := `
		INSERT INTO placement_coordinators (user_id, usn, branch, batch)
		VALUES ($1, $2, $3, $4)
	`
	_, err = tx.Exec(coordinatorQuery, adminID, coordinator.USN, coordinator.Branch, coordinator.Batch)
	if err != nil {
		log.Printf("Error inserting placement coordinator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Placement coordinator added successfully"))

}

// CreatePlacementCoordinatorsTable sets up table with foreign key constraints
// Links to both admins and students tables for data integrity
func CreatePlacementCoordinatorsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS placement_coordinators (
		user_id VARCHAR(50) NOT NULL REFERENCES admins(id) ON DELETE CASCADE, -- Foreign key linking to admins table
		usn VARCHAR(15) PRIMARY KEY REFERENCES students(usn),          -- USN of the placement coordinator
		branch VARCHAR(50) NOT NULL,               -- Branch of the placement coordinator
		batch VARCHAR(10) NOT NULL
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating placement_coordinators table: %v", err)
	} else {
		log.Println("Placement coordinators table ensured to exist.")
	}
}

// GetAllPlacementCoordinators retrieves coordinator details with joins
// Returns combined data from admins and placement_coordinators tables
func GetAllPlacementCoordinators(w http.ResponseWriter, r *http.Request, db *sql.DB) {

	if r.Method != http.MethodGet {
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
	if claims["role"] != "ADMIN" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	if db == nil {
		log.Println("Failed to initialize the database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Modify the query to include role
	query := `
	SELECT pc.usn, a.name, a.email, pc.branch, pc.batch, a.contact, a.created_at, a.role
	FROM placement_coordinators pc
	INNER JOIN admins a ON pc.user_id = a.id
	`
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Error fetching placement coordinators: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var coordinators []models.PlacementCoordinator
	for rows.Next() {
		var pc models.PlacementCoordinator
		if err := rows.Scan(
			&pc.USN,
			&pc.Name,
			&pc.Email,
			&pc.Branch,
			&pc.Batch,
			&pc.Contact,
			&pc.CreatedAt,
			&pc.Role, // Add role to scan
		); err != nil {
			log.Printf("Error scanning row: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		coordinators = append(coordinators, pc)
	}

	if len(coordinators) == 0 {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("No placement coordinators found"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(coordinators)
}

func DeletePlacementCoordinator(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
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
	if claims["role"] != "ADMIN" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	usn := r.URL.Query().Get("usn")
	if usn == "" {
		http.Error(w, "USN is required", http.StatusBadRequest)
		return
	}

	if db == nil {
		log.Println("Failed to initialize the database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	var userID string
	query := `SELECT user_id FROM placement_coordinators WHERE usn = $1`
	err = db.QueryRow(query, usn).Scan(&userID)
	if err == sql.ErrNoRows {
		http.Error(w, "Placement coordinator not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Error fetching placement coordinator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	deleteCoordinatorQuery := `DELETE FROM placement_coordinators WHERE usn = $1`
	_, err = tx.Exec(deleteCoordinatorQuery, usn)
	if err != nil {
		log.Printf("Error deleting placement coordinator: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	deleteAdminQuery := `DELETE FROM admins WHERE id = $1`
	_, err = tx.Exec(deleteAdminQuery, userID)
	if err != nil {
		log.Printf("Error deleting admin: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Placement coordinator deleted successfully"))
}

// EditPlacementCoordinator updates the details of a placement coordinator
func EditPlacementCoordinator(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPut {
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
	if claims["role"] != "ADMIN" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	var coordinator models.PlacementCoordinator
	if err = json.NewDecoder(r.Body).Decode(&coordinator); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if db == nil {
		log.Println("Failed to initialize the database")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	// Begin transaction
	tx, err := db.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Update `admins` table
	updateAdminQuery := `
		UPDATE admins
		SET name = $1, email = $2, contact = $3
		WHERE id = (SELECT user_id FROM placement_coordinators WHERE usn = $4)
	`
	_, err = tx.Exec(updateAdminQuery, coordinator.Name, coordinator.Email, coordinator.Contact, coordinator.USN)
	if err != nil {
		log.Printf("Error updating admin details: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Update `placement_coordinators` table
	updateCoordinatorQuery := `
		UPDATE placement_coordinators
		SET branch = $1, batch = $2
		WHERE usn = $3
	`
	_, err = tx.Exec(updateCoordinatorQuery, coordinator.Branch, coordinator.Batch, coordinator.USN)
	if err != nil {
		log.Printf("Error updating placement coordinator details: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Placement coordinator details updated successfully"))
}
