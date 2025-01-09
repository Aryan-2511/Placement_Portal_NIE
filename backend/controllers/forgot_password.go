package controllers
import(
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)
func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request,db *sql.DB,secretKey string) {
	var request struct {
		Email string `json:"email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	resetToken := utils.GenerateRandomString(32)
	query := `
		UPDATE students
		SET reset_token = $1
		WHERE college_email = $2
		RETURNING name
	`
	var name string
	err := db.QueryRow(query, resetToken, request.Email).Scan(&name)
	if err == sql.ErrNoRows {
		http.Error(w, "Email not found", http.StatusNotFound)
		return
	} else if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resetURL := "http://localhost:8080/reset-password?token=" + resetToken
	emailBody := "Hi " + name + ",<br><br>Click the link below to reset your password:<br>" +
		"<a href='" + resetURL + "'>Reset Password</a>"

	if err := utils.SendEmail(request.Email, "Password Reset", emailBody); err != nil {
		http.Error(w, "Error sending email", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password reset link sent to your email."))
}
