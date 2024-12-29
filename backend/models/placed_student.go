package models
import "time"

type PlacedStudent struct {
	ID            int       `json:"id"`
	USN           string    `json:"usn"`
	Name          string    `json:"name"`
	Email         string    `json:"email"`
	Branch        string    `json:"branch"`
	Company       string    `json:"company"`
	Package       float64   `json:"package"`
	PlacementDate time.Time `json:"placement_date"`
	Contact       string  	`json:"contact"`
}