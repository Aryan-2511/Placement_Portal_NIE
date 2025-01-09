package models

type User struct {
	Name                string  `json:"name"`
	USN                 string  `json:"usn"`
	DOB                 string  `json:"dob"`
	College_Email       string  `json:"college_email"`
	Personal_Email      string  `json:"personal_email"`
	Branch              string  `json:"branch"`
	Batch               string  `json:"batch"`
	Address             string  `json:"address"`
	Contact             string  `json:"contact"`
	Gender              string  `json:"gender"`
	Category            string  `json:"category"`
	Aadhar              string  `json:"aadhar"`
	PAN                 string  `json:"pan"`
	Class_10_Percentage float64 `json:"class_10_percentage"`
	Class_10_Year       int     `json:"class_10_year"`
	Class_10_Board      string  `json:"class_10_board"`
	Class_12_Percentage float64 `json:"class_12_percentage"`
	Class_12_Year       int     `json:"class_12_year"`
	Class_12_Board      string  `json:"class_12_board"`
	Current_CGPA        float64 `json:"current_cgpa"`
	Backlogs            int     `json:"backlogs"`
	Password            string  `json:"password"`
	Role                string  `json:"role"`
	IsPlaced            string  `json:"isPlaced"`
	Resume_link         string  `json:"resume_link"`
	IsVerified        	bool	`json:"is_verified"`
}
