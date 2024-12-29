package models

import (
	"time"
)

type Admin struct {
	Name 		string  `json:"name"`
	Password 	string  `json:"password"`
	Email    	string  `json:"email"`
	Contact  	string  `json:"contact"`
	Role     	string  `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}
