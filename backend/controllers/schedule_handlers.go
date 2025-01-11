package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "strings"
	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)
func AddEvent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR"{
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}


    var event models.Schedule
    if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Generate schedule ID based on batch
    var batchPart string
    if event.Batch != "" {
        batchPart = event.Batch[2:]
    } else {
        fmt.Print("Batch is empty")
    }

    queryCount := `SELECT COUNT(*) FROM schedule WHERE batch = $1`
    var count int
    err = db.QueryRow(queryCount, event.Batch).Scan(&count)
    if err != nil {
        http.Error(w, "Error generating schedule ID", http.StatusInternalServerError)
        return
    }

    scheduleID := fmt.Sprintf("SCH%s%03d", batchPart, count+1)
	tableName := "schedule"
	if utils.CheckTableExists(db, tableName) {
		fmt.Printf("Table '%s' exists.\n", tableName)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName)
		CreateScheduleTable(db)
	}
    query := `
        INSERT INTO schedule (schedule_id, title, description, start_time, end_time, created_by, batch)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
    _, err = db.Exec(query, scheduleID, event.Title, event.Description, event.StartTime, event.EndTime, event.CreatedBy, event.Batch)
    if err != nil {
        http.Error(w, "Error adding event", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    w.Write([]byte("Event added successfully"))
}
func CreateScheduleTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS schedule (
    schedule_id VARCHAR(50) PRIMARY KEY, -- Unique schedule ID
    title VARCHAR(255) NOT NULL,
    description TEXT,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    created_by VARCHAR(100), -- Admin or Coordinator
    batch VARCHAR(10), -- Batch for specific events
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating schedule table: %v", err)
	} else {
		log.Println("Schedule table ensured to exist.")
	}
}
func DeleteEvent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodDelete {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR"{
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
    // Extract schedule ID from the request
    scheduleID := r.URL.Query().Get("schedule_id")
    if scheduleID == "" {
        http.Error(w, "Schedule ID is required", http.StatusBadRequest)
        return
    }

    query := `DELETE FROM schedule WHERE schedule_id = $1`
    result, err := db.Exec(query, scheduleID)
    if err != nil {
        http.Error(w, "Error deleting event", http.StatusInternalServerError)
        log.Printf("Error deleting event: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Event not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Event deleted successfully"))
}
func EditEvent(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodPut {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR"{
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
    // Extract schedule ID from the request
    scheduleID := r.URL.Query().Get("schedule_id")
    if scheduleID == "" {
        http.Error(w, "Schedule ID is required", http.StatusBadRequest)
        return
    }

    // Decode the updated event details
    var updatedEvent models.Schedule
    if err := json.NewDecoder(r.Body).Decode(&updatedEvent); err != nil {
        http.Error(w, "Invalid input", http.StatusBadRequest)
        return
    }

    // Update query
    query := `
        UPDATE schedule
        SET title = $1, description = $2, start_time = $3, end_time = $4,
            created_by = $5, batch = $6
        WHERE schedule_id = $7
    `
    result, err := db.Exec(query,
        updatedEvent.Title,
        updatedEvent.Description,
        updatedEvent.StartTime,
        updatedEvent.EndTime,
        updatedEvent.CreatedBy,
        updatedEvent.Batch,
        scheduleID,
    )

    if err != nil {
        http.Error(w, "Error updating event", http.StatusInternalServerError)
        log.Printf("Error updating event: %v", err)
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        http.Error(w, "Event not found", http.StatusNotFound)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte("Event updated successfully"))
}
func GetAllEvents(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }
    // Extract the optional "batch" filter from query parameters
    batch := r.URL.Query().Get("batch")
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" && claims["role"]!="STUDENT"{
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

    var (
        query string
        rows  *sql.Rows
        err1   error
    )

    // Check if the batch filter is provided
    if batch != "" {
        query = "SELECT schedule_id, title, description, start_time, end_time, created_by, batch FROM schedule WHERE batch = $1 ORDER BY start_time"
        rows, err1 = db.Query(query, batch)
    } else {
        query = "SELECT schedule_id, title, description, start_time, end_time, created_by, batch FROM schedule ORDER BY start_time"
        rows, err1 = db.Query(query)
    }

    if err1 != nil {
        http.Error(w, "Error fetching events", http.StatusInternalServerError)
        log.Printf("Error executing query: %v", err1)
        return
    }
    defer rows.Close()

    var events []models.Schedule
    for rows.Next() {
        var event models.Schedule
        if err1 := rows.Scan(&event.ScheduleID, &event.Title, &event.Description, &event.StartTime, &event.EndTime, &event.CreatedBy, &event.Batch); err1 != nil {
            http.Error(w, "Error processing events", http.StatusInternalServerError)
            log.Printf("Error scanning row: %v", err1)
            return
        }
        events = append(events, event)
    }

    if rows.Err() != nil {
        http.Error(w, "Error iterating through rows", http.StatusInternalServerError)
        log.Printf("Error iterating rows: %v", rows.Err())
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)
}

func GetStudentEvents(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
		log.Print(err)
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" && claims["role"]!="STUDENT"{
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
    batch := r.URL.Query().Get("batch")
    if batch == "" {
        http.Error(w, "Batch is required", http.StatusBadRequest)
        return
    }

    query := "SELECT schedule_id, title, description, start_time, end_time FROM schedule WHERE batch = $1 ORDER BY start_time"
    rows, err := db.Query(query, batch)
    if err != nil {
        http.Error(w, "Error fetching events", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    var events []models.Schedule
    for rows.Next() {
        var event models.Schedule
        if err := rows.Scan(&event.ScheduleID, &event.Title, &event.Description, &event.StartTime, &event.EndTime); err != nil {
            http.Error(w, "Error processing events", http.StatusInternalServerError)
            return
        }
        events = append(events, event)
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(events)
}
