package controllers
import(
	"encoding/json"
	"log"
	"net/http"
	"Github.com/Aryan-2511/Placement_NIE/db"
)
func ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var request struct {
		Token       string `json:"token"`
		NewPassword string `json:"new_password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	hashedPassword, err := HashPassword(request.NewPassword)
	if err != nil {
		http.Error(w, "Error hashing password", http.StatusInternalServerError)
		return
	}

	query := `
		UPDATE students
		SET password_hash = $1, reset_token = NULL
		WHERE reset_token = $2
	`
	_, err = db.DB.Exec(query, hashedPassword, request.Token)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("Password reset successfully. You can now log in."))
}
