package controllers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
	"golang.org/x/crypto/bcrypt"
)


func SignupHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
    	log.Printf("Error reading request body: %v", err)
    	http.Error(w, "Invalid input", http.StatusBadRequest)
    	return	
	}
	log.Printf("Raw request body: %s", body)

	// Reset the body reader so it can be used again
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	// Parse JSON request body
	var user models.User
	err = json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if user.Password == "" {
		log.Println("Password is empty")
	}	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	fmt.Println("Hash:", string(hashedPassword))
    if err != nil {
        log.Fatal(err)
    }

	verificationToken := utils.GenerateRandomString(32) // Implement a function to generate random strings.

	// Insert into database
	tableName := "students"

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
	query := `
		INSERT INTO students (
			name, usn, dob, college_email, personal_email,branch, batch, address, contact, gender, category, aadhar, pan,
			class_10_percentage, class_10_year, class_10_board, class_12_percentage, class_12_year,
			class_12_board, current_cgpa, backlogs,password_hash, role, is_verified, verification_token,resume_link
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24, $25, $26)
	`
	_, err = db.Exec(query, user.Name, user.USN, user.DOB, user.College_Email, user.Personal_Email, user.Branch, user.Batch, user.Address, user.Contact, user.Gender, user.Category,
		user.Aadhar, user.PAN, user.Class_10_Percentage, user.Class_10_Year, user.Class_10_Board,
		user.Class_12_Percentage, user.Class_12_Year,user.Class_12_Board, user.Current_CGPA,
		user.Backlogs, hashedPassword, user.Role, false, verificationToken,user.Resume_link)
	
	if err != nil {
		log.Printf("Error inserting data: %v\n", err)
		http.Error(w, "Error inserting data", http.StatusInternalServerError)
		return
	}
	verificationURL := "http://localhost:8080/verify-email?token=" + verificationToken // Replace with your domain
	emailBody := "Welcome, " + user.Name + "!<br><br>Please verify your email by clicking the link below:<br>" +
		"<a href='" + verificationURL + "'>Verify Email</a>"

	if err := utils.SendEmail(user.College_Email, "Email Verification", emailBody); err != nil {
		http.Error(w, "Error sending verification email", http.StatusInternalServerError)
		return
	}
	// Respond to client
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User  registered successfully"})
}

func CreateStudentsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS students (
			name VARCHAR(100) NOT NULL,
			usn VARCHAR(10) PRIMARY KEY,
			dob DATE NOT NULL,
			college_email VARCHAR(100) UNIQUE NOT NULL,
			personal_email VARCHAR(100) UNIQUE NOT NULL,
			branch VARCHAR(50),
			batch VARCHAR(10),
			address TEXT,
			contact VARCHAR(15),
			gender VARCHAR(10),
			category VARCHAR(50),
			aadhar VARCHAR(16),
			pan VARCHAR(15),
			class_10_percentage DECIMAL(5, 2),
			class_10_year INT,
			class_10_board VARCHAR(50),
			class_12_percentage DECIMAL(5, 2),
			class_12_year INT,
			class_12_board VARCHAR(50),
			current_cgpa DECIMAL(4, 2),
			backlogs INT DEFAULT 0,
			password_hash VARCHAR(100) NOT NULL,
			role VARCHAR(20) NOT NULL DEFAULT 'STUDENT',
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			is_verified BOOLEAN DEFAULT FALSE,
			verification_token TEXT,
			reset_token TEXT,
			isPlaced VARCHAR(15) DEFAULT 'NO',
			resume_link VARCHAR(100)
		);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	} else {
		log.Println("Table `students` created or already exists.")
	}
}