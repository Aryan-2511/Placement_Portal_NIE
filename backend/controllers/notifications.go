package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	// "fmt"
	// "log"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

// NotificationRequest defines the structure for notification requests
type NotificationRequest struct {
	Title            string   `json:"title"`
	Message          string   `json:"message"`
	NotificationType string   `json:"notification_type"`
	TargetType       string   `json:"target_type"`
	TargetValue      string   `json:"target_value"`
	CustomEmails     []string `json:"custom_emails,omitempty"`
}

// // GenerateNotificationID generates a unique notification ID
// func GenerateNotificationID(year string, serial int) string {
//     return fmt.Sprintf("NOT%s%04d", year[2:], serial)
// }

// Add this function to get recipients
func getRecipients(db *sql.DB, targetType, targetValue string, customEmails []string) ([]string, error) {
	var recipients []string

	switch targetType {
	case "BATCH":
		rows, err := db.Query("SELECT college_email FROM students WHERE batch = $1", targetValue)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				return nil, err
			}
			recipients = append(recipients, email)
		}
	case "CUSTOM":
		recipients = customEmails
	case "ALL":
		rows, err := db.Query("SELECT college_email FROM students")
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var email string
			if err := rows.Scan(&email); err != nil {
				return nil, err
			}
			recipients = append(recipients, email)
		}
	}

	return recipients, nil
}

func SendUnifiedNotification(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	// Add auth check
	claims, err := utils.ValidateToken(r.Header.Get("Authorization"))
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req NotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Generate notification ID
	notificationID := GenerateNotificationID(time.Now().Format("2006"), 1) // You'll need to implement proper serial number generation

	// Get recipient emails based on target type
	recipients, err := getRecipients(db, req.TargetType, req.TargetValue, req.CustomEmails)
	if err != nil {
		http.Error(w, "Failed to get recipients", http.StatusInternalServerError)
		return
	}

	// Insert notification
	tx, err := db.Begin()
	if err != nil {
		http.Error(w, "Transaction error", http.StatusInternalServerError)
		return
	}

	// Insert main notification
	_, err = tx.Exec(`
        INSERT INTO notifications (id, title, message, notification_type, target_type, target_value, created_by)
        VALUES ($1, $2, $3, $4, $5, $6, $7)`,
		notificationID, req.Title, req.Message, req.NotificationType, req.TargetType, req.TargetValue, claims["id"])

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to create notification", http.StatusInternalServerError)
		return
	}

	// Insert recipients
	for _, email := range recipients {
		_, err = tx.Exec(`
            INSERT INTO notification_recipients (notification_id, recipient_email)
            VALUES ($1, $2)`,
			notificationID, email)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to add recipients", http.StatusInternalServerError)
			return
		}
	}

	if err = tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	// Send emails if required
	if req.NotificationType == "EMAIL" || req.NotificationType == "BOTH" {
		go utils.SendEmailsAsync(recipients, req.Title, req.Message, notificationID, db)
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification created successfully",
		"id":      notificationID,
	})
}
