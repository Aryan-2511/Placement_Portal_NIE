package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
	"golang.org/x/crypto/bcrypt"
)
func GenerateAdminID(role string, serial int) string {
	roleCode := map[string]string{
		"ADMIN":               "PO", // Default role mapping
		"PLACEMENT_COORDINATOR": "PC",
	}

	// Format serial number as a 3-digit string
	serialStr := fmt.Sprintf("%03d", serial)

	// Return formatted Admin-ID
	return fmt.Sprintf("AD%s%s", roleCode[role], serialStr)
}

func AddAdmin(w http.ResponseWriter,r *http.Request,db *sql.DB){
	if r.Method != http.MethodPost{
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
	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err!=nil{
		http.Error(w,"Invalid request payload",http.StatusBadRequest)
		return
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error hashing password: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	tableName := "admins"

	exists, err := utils.CheckTableExists(db, tableName)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		return
	}

	if exists {
		fmt.Printf("Table '%s' exists.\n", tableName)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName)
		CreateApplicationsTable(db)
	}
	var serial int
	query := `SELECT COUNT(*) + 1 FROM admins`
	err = db.QueryRow(query).Scan(&serial)
	if err != nil {
		log.Printf("Error fetching serial for Admin-ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	adminID := GenerateAdminID(admin.Role, serial)
	query = `
			INSERT INTO admins(id, name, password_hash, email, contact, role, created_at)
			VALUES($1, $2, $3, $4, $5, $6, NOW())
			`
	_,err = db.Exec(query,adminID,admin.Name, hashedPassword, admin.Email, admin.Contact, admin.Role)
	if err != nil{
		log.Printf("Error adding admin: %v", err)
		http.Error(w,"Internal server error", http.StatusInternalServerError)
		return 
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Admin added successfully"))

}
func CreateAdminsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS admins (
		id VARCHAR(50) PRIMARY KEY,                     -- Unique ID for each admin
		name VARCHAR(255) NOT NULL,               -- Name of the admin
		password_hash TEXT NOT NULL,              -- Hashed password
		email VARCHAR(255) UNIQUE NOT NULL,       -- Unique email address
		contact VARCHAR(15) NOT NULL,             -- Contact number
		role VARCHAR(50) NOT NULL,                -- Role (e.g., "admin", "placement_coordinator")
		created_at TIMESTAMP DEFAULT NOW()        -- Timestamp for when the admin was added
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating admins table: %v", err)
	} else {
		log.Println("Admins table ensured to exist.")
	}
}
func EditAdmin(w http.ResponseWriter, r *http.Request,db *sql.DB){
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

	var admin models.Admin
	if err := json.NewDecoder(r.Body).Decode(&admin); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if admin.Email == "" {
		http.Error(w, "Email is required to identify the admin", http.StatusBadRequest)
		return
	}
	
	query := `
			UPDATE admins
			SET name = $1 , contact = $2
			WHERE email = $3
	`
	_,err = db.Exec(query,
		admin.Name,
		admin.Contact,
		admin.Email,
	)
	if err != nil {
		log.Printf("Error updating admin details: %v", err)
		http.Error(w, "Error updating admin details", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Admin details updated successfully"))	

}