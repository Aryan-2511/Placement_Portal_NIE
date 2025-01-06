package models

import "time"

type Schedule struct {
	ScheduleID  string    `json:"schedule_id"`            // Unique schedule ID
	Title       string    `json:"title"`                  // Title of the event
	Description string    `json:"description,omitempty"`  // Description of the event
	StartTime   time.Time `json:"start_time"`             // Start time of the event
	EndTime     time.Time `json:"end_time"`               // End time of the event
	CreatedBy   string    `json:"created_by,omitempty"`   // Created by (Admin or Coordinator)
	Batch       string    `json:"batch,omitempty"`        // Batch (if applicable for specific events)
	CreatedAt   time.Time `json:"created_at"`             // Timestamp when the event was created
}