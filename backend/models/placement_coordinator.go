package models
import "time"

type PlacementCoordinator struct {
	USN    		string `json:"usn"`
	Branch 		string `json:"branch"`
	Batch  		string `json:"batch"`
	Name 		string  `json:"name"`
	Password 	string  `json:"password"`
	Email    	string  `json:"email"`
	Contact  	string  `json:"contact"`
	Role      string `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}
