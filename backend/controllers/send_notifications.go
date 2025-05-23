package controllers
import(
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/utils"
)
// SendNotificationHandler manages bulk email notifications
// Supports batch, opportunity, and custom recipient targeting
func SendNotificationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Parse request payload
	var payload struct {
		Criteria    string   `json:"criteria"`    // "batch", "opportunity", "custom"
		Value       string   `json:"value"`       // Batch year, Opportunity ID, or unused for custom
		CustomEmails []string `json:"customEmails"` // Used when Criteria is "custom"
		Subject     string   `json:"subject"`
		Message     string   `json:"message"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var emails []string
	switch payload.Criteria {
	case "batch":
		emails, err = utils.GetEmailsByBatch(db, payload.Value)
	case "opportunity":
		emails, err = utils.GetEmailsByOpportunity(db, payload.Value)
	case "custom":
		emails = payload.CustomEmails
	default:
		http.Error(w, "Invalid criteria", http.StatusBadRequest)
		return
	}

	if err != nil {
		http.Error(w, "Failed to fetch email list", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	// Send emails in bulk
	if err := utils.SendBulkEmail(emails, payload.Subject, payload.Message); err != nil {
		http.Error(w, "Failed to send notifications", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notifications sent successfully"))
}
