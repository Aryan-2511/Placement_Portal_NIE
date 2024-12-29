package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"Github.com/Aryan-2511/Placement_NIE/db"
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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var hashedPassword, name string
	var verificationStatus bool

	var query string
	switch loginRequest.Role {
	case "ADMIN":
		query = "SELECT password_hash, name FROM admins WHERE email = $1"
	case "PLACEMENT_COORDINATOR":
		query = `
			SELECT pc.password_hash, a.name
			FROM placement_coordinators pc
			JOIN admins a ON pc.user_id = a.id
			WHERE a.email = $1`
	case "STUDENT":
		query = "SELECT password_hash, name, is_verified FROM students WHERE college_email = $1"
	default:
		http.Error(w, "Invalid role", http.StatusBadRequest)
		return
	}

	err := db.DB.QueryRow(query, loginRequest.Email).Scan(&hashedPassword, &name, &verificationStatus)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if loginRequest.Role == "STUDENT" && !verificationStatus {
		http.Error(w, "Please verify your email before logging in", http.StatusForbidden)
		return
	}

	if err := CheckPasswordHash(loginRequest.Password, hashedPassword); err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		fmt.Print(err)
		return
	}

	token, err := utils.GenerateToken(loginRequest.Email, loginRequest.Role)
	if err != nil {
		log.Printf("Token generation error: %v", err)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(map[string]string{
		"token": token,
		"role":  loginRequest.Role,
		"name":  name,
	}); err != nil {
		log.Printf("Error encoding response: %v", err)
		http.Error(w, "Error generating response", http.StatusInternalServerError)
		return
	}
}