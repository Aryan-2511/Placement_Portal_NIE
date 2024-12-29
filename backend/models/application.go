package models

import "time"

// Application represents an application record
type Application struct {
	ID            	int       	`json:"id,omitempty"`
	StudentUSN    	string    	`json:"student_usn"`
	StudentName   	string    	`json:"student_name"`
	OpportunityID 	string      `json:"opportunity_id"`
	AppliedAt     	time.Time 	`json:"applied_at"`
	Status        	string 		`json:"status"`
}
