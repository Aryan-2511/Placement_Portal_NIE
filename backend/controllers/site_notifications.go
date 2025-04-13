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
	// Take last two digits of the batch year (e.g., "25" from "2025")
	batchCode := batch[2:4] // Get "25" from "2025"
	serialStr := fmt.Sprintf("%04d", serial)
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
		StudentEmails []string `json:"student_emails"`
		Title         string   `json:"title"`
		Message       string   `json:"message"`
		Batch         string   `json:"batch"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the next serial number for notification ID
	var serial int
	err = db.QueryRow(`
	    SELECT COALESCE(MAX(CAST(SUBSTRING(id, 5) AS INTEGER)), 0) + 1 
	    FROM notifications 
	    WHERE batch = $1
	`, payload.Batch).Scan(&serial)
	if err != nil {
		log.Printf("Error getting serial: %v", err)
		http.Error(w, "Failed to generate notification ID", http.StatusInternalServerError)
		return
	}

	// Generate notification ID
	notificationID := GenerateNotificationID(payload.Batch, serial)

	// // Get admin ID from claims
	// adminID, ok := claims["admin_id"].(string)
	// if !ok {
	// 	log.Printf("Error: admin_id not found in token claims")
	// 	http.Error(w, "Invalid token claims", http.StatusBadRequest)
	// 	return
	// }

	// Insert the notification into the notifications table
	query := `
		INSERT INTO notifications (
			id, title, message, notification_type, batch, 
			target_type, target_value, created_by, created_at
		) VALUES (
			$1, $2, $3, $4, $5, 
			$6, $7, $8, NOW()
		)
	`
	_, err = db.Exec(
		query,
		notificationID,
		payload.Title,
		payload.Message,
		"SITE",
		payload.Batch,
		"BATCH",
		payload.Batch,
		claims["role"], // Using the admin ID instead of role
	)
	if err != nil {
		log.Printf("Error creating notification: %v", err)
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	// Insert into notification_recipients
	// In AddNotificationHandler, modify the INSERT query for notification_recipients:
	if len(payload.StudentEmails) == 0 {
		// Global notification for batch
		query = `
	        INSERT INTO notification_recipients (notification_id, recipient_email, is_sent, sent_at)
	        SELECT $1, college_email, TRUE, NOW() FROM students WHERE batch = $2
	    `
		_, err = db.Exec(query, notificationID, payload.Batch)
		if err != nil {
			log.Printf("Error mapping global notification: %v", err)
			http.Error(w, "Failed to map global notification", http.StatusInternalServerError)
			return
		}
	} else {
		// Specific students
		for _, email := range payload.StudentEmails {
			_, err = db.Exec(`
	            INSERT INTO notification_recipients (
	                notification_id, 
	                recipient_email, 
	                is_sent, 
	                sent_at
	            )
	            VALUES ($1, $2, TRUE, NOW())
	        `, notificationID, email)
			if err != nil {
				log.Printf("Error mapping notification to student %s: %v", email, err)
				http.Error(w, fmt.Sprintf("Failed to map notification to student %s", email), http.StatusInternalServerError)
				return
			}
		}
	}

	// In MarkNotificationAsReadHandler, modify the UPDATE query:
	query = `
	    UPDATE notification_recipients
	    SET is_read = TRUE,
	        read_at = NOW()
	    WHERE notification_id = $1 AND recipient_email = $2
	`
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification created successfully",
		"id":      notificationID,
	})
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
        created_at TIMESTAMP DEFAULT NOW()
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
	studentEmail, ok := claims["email"].(string)
	if !ok {
		log.Printf("Error: email not found in token claims: %v", claims)
		http.Error(w, "Invalid token claims", http.StatusBadRequest)
		return
	}
	fmt.Print("Student Email: ", studentEmail)

	// Fetch notifications for the student
	query := `
        SELECT DISTINCT n.id, n.title, n.message, 
            COALESCE(nr.is_read, FALSE) as is_read, 
            n.created_at
        FROM notifications n
        LEFT JOIN notification_recipients nr ON n.id = nr.notification_id 
            AND nr.recipient_email = $1
        WHERE n.batch = (
            SELECT batch FROM students WHERE college_email = $1
        )
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
			ID        string `json:"id"` // Changed from int to string
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
	studentEmail := claims["email"]

	// Parse notification ID from request body
	var payload struct {
		NotificationID string `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Mark notification as read
	query := `
		UPDATE notification_recipients
		SET is_read = TRUE
		WHERE notification_id = $1 AND recipient_email = $2
	`
	if _, err := db.Exec(query, payload.NotificationID, studentEmail); err != nil {
		http.Error(w, "Failed to update notification", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Notification marked as read"))
}
