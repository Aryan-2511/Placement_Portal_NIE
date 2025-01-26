package models

import "time"


type Branch string

// List of valid branches as constants
const (
	BranchCSE   Branch = "CSE"
	BranchISE   Branch = "ISE"
	BranchECE   Branch = "ECE"
	BranchEEE   Branch = "EEE"
	BranchMECH  Branch = "MECH"
	BranchCIVIL Branch = "CIVIL"
	BranchAIML  Branch = "AIML"
)

type Coordinator struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type Opportunity struct {
	ID                				string    `json:"id"`
	Title             				string    `json:"title"`
	Company           				string    `json:"company"`
	Batch             				string    `json:"batch"`
	Location          				string    `json:"location"`
	CTC               				float64   `json:"ctc"`
	CTCInfo           				string    `json:"ctc_info"`
	CGPA              				float64   `json:"cgpa"`
	Category          				string    `json:"category"`
	Backlog           				int       `json:"backlog"`
	AllowedBranches   				[]Branch  `json:"allowed_branches"`
	AllowedGenders    				[]string  `json:"allowed_genders"` 
	RegistrationDate  				time.Time `json:"registration_date"`
	Coordinators      				[]Coordinator `json:"coordinators"`
	AdditionalInfo    				string    `json:"additional_info"`
	FormLink          				string    `json:"form_link"`
	JobDescription    				string    `json:"job_description"`
	AttachedDocuments 				[]string  `json:"attached_documents"`
	OpportunityType   				string    `json:"opportunity_type"`
	Created_By						string 	  `json:"created_by"`
    Class_10_Percentage_Criteria 	float64   `json:"class_10_percentage_criteria"`
	Class_12_Percentage_Criteria 	float64   `json:"class_12_percentage_criteria"`
	Status 							string 	  `json:"status"`
	Completed 						string    `json:"completed"`
	CreatedAt						time.Time `json:"created_at"`
}
