package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
	"golang.org/x/crypto/bcrypt"
)
func CheckPasswordHash(password, hash string) error {
	log.Printf("Checking password: '%s' against hash: '%s'", password, hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        fmt.Println("Password does not match!")
    } else {
        fmt.Println("Password matches!")
    }
	return err
}

func LoginHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse login request payload
	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var (
		hashedPassword string
		userDetails    map[string]interface{}
		query          string
	)

	// // Define role-specific structs
	// type Student struct {
	// 	Name              string
	// 	USN               string
	// 	DOB               string
	// 	CollegeEmail      string
	// 	PersonalEmail     string
	// 	Branch            string
	// 	Batch             string
	// 	Address           string
	// 	Contact           string
	// 	Gender            string
	// 	Category          string
	// 	Aadhar            string
	// 	PAN               string
	// 	Class10Percentage float64
	// 	Class10Year       int
	// 	Class10Board      string
	// 	Class12Percentage float64
	// 	Class12Year       int
	// 	Class12Board      string
	// 	CurrentCGPA       float64
	// 	Backlogs          int
	// 	IsVerified        bool
	// 	ResumeLink        string
	// }

	// type Admin struct {
	// 	ID      int
	// 	Name    string
	// 	Email   string
	// 	Contact string
	// }

	// Role-based query selection
	switch loginRequest.Role {
	case "ADMIN":
		query = "SELECT password_hash, id, name, email, contact FROM admins WHERE email = $1"
	case "PLACEMENT_COORDINATOR":
		query = `
			SELECT a.password_hash, a.id, a.name, a.email, pc.contact
			FROM placement_coordinators pc
			JOIN admins a ON pc.user_id = a.id
			WHERE a.email = $1`
	case "STUDENT":
		query = `
			SELECT password_hash, name, usn, dob, college_email, personal_email, branch, batch, address, contact, gender, category, 
			       aadhar, pan, class_10_percentage, class_10_year, class_10_board, class_12_percentage, 
			       class_12_year, class_12_board, current_cgpa, backlogs, isPlaced, resume_link, is_verified
			FROM students 
			WHERE college_email = $1`
	default:
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	// Fetch data based on role
	userDetails = make(map[string]interface{})
	switch loginRequest.Role {
	case "ADMIN", "PLACEMENT_COORDINATOR":
		var admin models.Admin
		err := db.QueryRow(query, loginRequest.Email).Scan(&hashedPassword, &admin.ID, &admin.Name, &admin.Email, &admin.Contact)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Populate userDetails
		userDetails["id"] = admin.ID
		userDetails["name"] = admin.Name
		userDetails["email"] = admin.Email
		userDetails["contact"] = admin.Contact

	case "STUDENT":
		var student models.User
		err := db.QueryRow(query, loginRequest.Email).Scan(
			&hashedPassword, &student.Name, &student.USN, &student.DOB, &student.College_Email, &student.Personal_Email,
			&student.Branch, &student.Batch, &student.Address, &student.Contact, &student.Gender, &student.Category,
			&student.Aadhar, &student.PAN, &student.Class_10_Percentage, &student.Class_10_Year, &student.Class_10_Board,
			&student.Class_12_Percentage, &student.Class_12_Year, &student.Class_12_Board, &student.Current_CGPA,
			&student.Backlogs,&student.IsPlaced, &student.Resume_link, &student.IsVerified,
		)
		if err == sql.ErrNoRows {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			return
		} else if err != nil {
			log.Printf("Database error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		// Check verification status
		if !student.IsVerified {
			http.Error(w, "Please verify your email before logging in", http.StatusForbidden)
			return
		}

		// Populate userDetails
		userDetails = map[string]interface{}{
			"name":          student.Name,
			"usn":           student.USN,
			"dob":           student.DOB,
			"college_email": student.College_Email,
			"personal_email": student.Personal_Email,
			"branch": student.Branch,
			"batch": student.Batch,
			"address": student.Address,
			"contact": student.Contact,
			"gender": student.Gender,
			"category": student.Category,
			"aadhar": student.Aadhar,
			"pan": student.PAN,
			"class_10_percentage": student.Class_10_Percentage,
			"class_10_year": student.Class_10_Year,
			"class_10_board": student.Class_10_Board,
			"class_12_percentage": student.Class_12_Percentage,
			"class_12_year": student.Class_12_Year,
			"class_12_board": student.Class_12_Board,
			"current_cgpa": student.Current_CGPA,
			"backlogs": student.Backlogs,
			"isPlaced": student.IsPlaced,
			"resume_link": student.Resume_link,
		}
	}

	// Validate password
	if err := CheckPasswordHash(loginRequest.Password, hashedPassword); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT token
	token, err := utils.GenerateToken(loginRequest.Email, loginRequest.Role, userDetails["name"].(string))
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	// Include token in user details
	userDetails["token"] = token
	userDetails["role"] = loginRequest.Role

	// Send response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(userDetails); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}
