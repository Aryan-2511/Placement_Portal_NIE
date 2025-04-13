package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/utils"
)

func GetBatchNotificationsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Get batch from query parameters
	batch := r.URL.Query().Get("batch")
	if batch == "" {
		http.Error(w, "Batch parameter is required", http.StatusBadRequest)
		return
	}

	// Fetch notifications for the batch
	query := `
        SELECT n.id, n.title, n.message, n.created_at,
            COUNT(nr.recipient_email) as total_recipients,
            COUNT(CASE WHEN nr.is_read = true THEN 1 END) as read_count
        FROM notifications n
        LEFT JOIN notification_recipients nr ON n.id = nr.notification_id
        WHERE n.batch = $1
        GROUP BY n.id, n.title, n.message, n.created_at
        ORDER BY n.created_at DESC
    `
	rows, err := db.Query(query, batch)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		log.Print(err)
		return
	}
	defer rows.Close()

	var notifications []map[string]interface{}
	for rows.Next() {
		var notification struct {
			ID              string `json:"id"`
			Title           string `json:"title"`
			Message         string `json:"message"`
			CreatedAt       string `json:"created_at"`
			TotalRecipients int    `json:"total_recipients"`
			ReadCount       int    `json:"read_count"`
		}
		if err := rows.Scan(
			&notification.ID,
			&notification.Title,
			&notification.Message,
			&notification.CreatedAt,
			&notification.TotalRecipients,
			&notification.ReadCount,
		); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			log.Print(err)
			return
		}
		notifications = append(notifications, map[string]interface{}{
			"id":               notification.ID,
			"title":            notification.Title,
			"message":          notification.Message,
			"created_at":       notification.CreatedAt,
			"total_recipients": notification.TotalRecipients,
			"read_count":       notification.ReadCount,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(notifications)
}

func DeleteNotificationHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

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

	// Parse notification ID from request body
	var payload struct {
		NotificationID string `json:"notification_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Delete notification (cascade will handle recipients)
	query := `DELETE FROM notifications WHERE id = $1`
	result, err := db.Exec(query, payload.NotificationID)
	if err != nil {
		http.Error(w, "Failed to delete notification", http.StatusInternalServerError)
		log.Print(err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Error checking deletion result", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Notification not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Notification deleted successfully",
	})
}
