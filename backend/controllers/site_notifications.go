package controllers
import(
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"fmt"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)
func GenerateNotificationID(batch string, serial int) string {
	serialStr := fmt.Sprintf("%04d", serial) // 3-digit serial
	batchCode := batch[len(batch)-2:]        // Last two digits of the batch
	return fmt.Sprintf("NT%s%s", batchCode, serialStr)
}

func AddNotificationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	tableName := "notifications"

	exists, err := utils.CheckTableExists(db, tableName)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		return
	}

	if exists {
		fmt.Printf("Table '%s' exists.\n", tableName)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName)
		CreateNotificationsTable(db)
	}

	// Parse the request payload
	var payload struct {
		StudentIDs []string  `json:"student_emails"` // Empty for global notification
		Title      string `json:"title"`
		Message    string `json:"message"`
		Batch      string `json:"batch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	var serial int
	query := `SELECT COUNT(*) + 1 FROM opportunities WHERE batch = $1`
	err = db.QueryRow(query, payload.Batch).Scan(&serial)
	if err != nil {
		log.Printf("Error fetching serial for Opportunity-ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	

	// Insert notifications into the database
	query = `
		INSERT INTO notifications (student_emails, title,batch, message)
		VALUES ($1, $2, $3,$4)
	`
	for _, studentID := range payload.StudentIDs {
		if _, err := db.Exec(query, studentID, payload.Title, payload.Message); err != nil {
			log.Printf("Error inserting notification for student %s: %v", studentID, err)
		}
	}

	// Broadcast notification (global)
	if len(payload.StudentIDs) == 0 {
		query = `
			INSERT INTO notifications (student_id, title, message)
			SELECT id, $1, $2 FROM students
		`
		if _, err := db.Exec(query, payload.Title, payload.Message); err != nil {
			log.Printf("Error inserting global notifications: %v", err)
			http.Error(w, "Failed to create notifications", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Notifications created successfully"))
}
func CreateNotificationsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS notifications (
    id SERIAL PRIMARY KEY,
    student_email TEXT, -- Null if it's a global notification
    title TEXT NOT NULL,
	batch TEXT NOT NULL,
    message TEXT NOT NULL,
    is_read BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT NOW()
	);

	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating notifications table: %v", err)
	} else {
		log.Println("Notifications table ensured to exist.")
	}
}
func GetNotificationsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if claims["role"] != "STUDENT" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Extract student ID from token claims
	studentEmail := claims["college_email"]

	// Fetch notifications for the student
	query := `
		SELECT id, title, message, is_read, created_at
		FROM notifications
		WHERE student_email = $1 OR student_email IS NULL
		ORDER BY created_at DESC
	`
	rows, err := db.Query(query, studentEmail)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var notification struct {
			ID        int    `json:"id"`
			Title     string `json:"title"`
			Message   string `json:"message"`
			IsRead    bool   `json:"is_read"`
			CreatedAt string `json:"created_at"`
		}
		if err := rows.Scan(&notification.ID, &notification.Title, &notification.Message, &notification.IsRead, &notification.CreatedAt); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Print(err)
			return
		}
		notifications = append(notifications, map[string]interface{}{
			"id":         notification.ID,
			"title":      notification.Title,
			"message":    notification.Message,
			"is_read":    notification.IsRead,
			"created_at": notification.CreatedAt,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func MarkNotificationAsReadHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if claims["role"] != "STUDENT" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Parse notification ID
	var payload struct {
		NotificationID int `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mark notification as read
	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1
	`
	if _, err := db.Exec(query, payload.NotificationID); err != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification marked as read"))
}
