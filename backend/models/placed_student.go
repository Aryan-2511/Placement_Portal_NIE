package models

import "time"

type PlacedStudent struct {
	ID            string       `json:"id"`
	USN           string    `json:"usn"`
	OpportunityID string    `json:"opportunity_id"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Branch        string    `json:"branch"`
	Batch         string    `json:"batch"`
	Company       string    `json:"company"`
	Package       float64   `json:"package"`
	PlacementDate time.Time `json:"placement_date"`
	Contact       string  	`json:"contact"`
	PlacementType string    `json:"placement_type"`
}