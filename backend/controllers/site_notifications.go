package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Parse the request payload
	var payload struct {
		StudentEmails []string `json:"student_emails"` // Empty for global notifications
		Title         string   `json:"title"`
		Message       string   `json:"message"`
		Batch         string   `json:"batch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Insert the notification into the notifications table
	query := `
		INSERT INTO notifications (title, batch, message, created_by)
		VALUES ($1, $2, $3, $4) RETURNING id
	`
	var notificationID int
	err = db.QueryRow(query, payload.Title, payload.Batch, payload.Message, claims["id"]).Scan(&notificationID)
	if err != nil {
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	// If it's a global notification, map to all students of the batch
	if len(payload.StudentEmails) == 0 {
		query = `
			INSERT INTO notification_students (notification_id, student_email)
			SELECT $1, email FROM students WHERE batch = $2
		`
		_, err = db.Exec(query, notificationID, payload.Batch)
		if err != nil {
			http.Error(w, "Failed to map global notification", http.StatusInternalServerError)
			return
		}
	} else {
		// Map the notification to specific students
		for _, email := range payload.StudentEmails {
			_, err = db.Exec(`
				INSERT INTO notification_students (notification_id, student_email)
				VALUES ($1, $2)
			`, notificationID, email)
			if err != nil {
				http.Error(w, fmt.Sprintf("Failed to map notification to student %s", email), http.StatusInternalServerError)
				return
			}
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Notification created successfully"))
}

func CreateNotificationsTable(db *sql.DB) {
    query := `
    CREATE TABLE IF NOT EXISTS notifications (
        id TEXT PRIMARY KEY,
        title TEXT NOT NULL,
        message TEXT NOT NULL,
        notification_type TEXT NOT NULL, -- 'EMAIL', 'SITE', 'BOTH'
        batch TEXT,
        target_type TEXT NOT NULL, -- 'BATCH', 'OPPORTUNITY', 'CUSTOM', 'ALL'
        target_value TEXT, -- batch year, opportunity_id, or null for custom/all
        created_by TEXT NOT NULL,
        created_at TIMESTAMP DEFAULT NOW(),
        FOREIGN KEY (created_by) REFERENCES admins(id)
    );

    CREATE TABLE IF NOT EXISTS notification_recipients (
        notification_id TEXT NOT NULL,
        recipient_email TEXT NOT NULL,
        is_read BOOLEAN DEFAULT FALSE,
        is_sent BOOLEAN DEFAULT FALSE,
        sent_at TIMESTAMP,
        read_at TIMESTAMP,
        PRIMARY KEY (notification_id, recipient_email),
        FOREIGN KEY (notification_id) REFERENCES notifications(id) ON DELETE CASCADE
    );`

    _, err := db.Exec(query)
    if err != nil {
        log.Printf("Error creating notifications tables: %v", err)
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

	// Extract student email from token claims
	studentEmail := claims["college_email"]

	// Fetch notifications for the student
	query := `
		SELECT n.id, n.title, n.message, ns.is_read, n.created_at
		FROM notifications n
		JOIN notification_students ns ON n.id = ns.notification_id
		WHERE ns.student_email = $1
		ORDER BY n.created_at DESC
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

	// Extract student email from token claims
	studentEmail := claims["college_email"]

	// Parse notification ID from request body
	var payload struct {
		NotificationID int `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mark notification as read
	query := `
		UPDATE notification_students
		SET is_read = TRUE
		WHERE notification_id = $1 AND student_email = $2
	`
	if _, err := db.Exec(query, payload.NotificationID, studentEmail); err != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification marked as read"))
}
