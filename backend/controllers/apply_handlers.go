package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

func ApplyHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var request struct {
		StudentUSN    string `json:"student_usn"`
		OpportunityID string    `json:"opportunity_id"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	var opportunityStatus string
	err = db.QueryRow("SELECT status FROM opportunities WHERE id = $1", request.OpportunityID).Scan(&opportunityStatus)
	if err != nil {
		http.Error(w, "Opportunity not found", http.StatusNotFound)
		return
	}

	if opportunityStatus == "Closed" {
		http.Error(w, "Opportunity is closed for applications", http.StatusForbidden)
		return
	}

	// Fetch necessary student details
	var student models.User
	studentQuery := `SELECT usn, name, current_cgpa, class_10_percentage, class_12_percentage, branch, batch, backlogs, gender 
                     FROM students WHERE usn = $1`
	err = db.QueryRow(studentQuery, request.StudentUSN).Scan(
		&student.USN,
		&student.Name,
		&student.Current_CGPA,
		&student.Class_10_Percentage,
		&student.Class_12_Percentage,
		&student.Branch,
		&student.Batch,
		&student.Backlogs,
		&student.Gender,
	)
	if err != nil {
		http.Error(w, "Student not found", http.StatusNotFound)
		log.Print(err)
		return
	}

	// Fetch necessary opportunity details
	var allowedBranchesRaw, allowedGendersRaw []byte
	var opportunity models.Opportunity
	opportunityQuery := `SELECT id, cgpa, class_10_percentage_criteria, class_12_percentage_criteria, allowed_branches, 
                         batch, backlog, allowed_genders, registration_date 
                         FROM opportunities WHERE id = $1`
	err = db.QueryRow(opportunityQuery, request.OpportunityID).Scan(
		&opportunity.ID,
		&opportunity.CGPA,
		&opportunity.Class_10_Percentage_Criteria,
		&opportunity.Class_12_Percentage_Criteria,
		&allowedBranchesRaw,
		&opportunity.Batch,
		&opportunity.Backlog,
		&allowedGendersRaw,
		&opportunity.RegistrationDate,
	)
	if err != nil {
		http.Error(w, "Opportunity not found", http.StatusNotFound)
		log.Print(err)
		return
	}
	if err := json.Unmarshal(allowedBranchesRaw, &opportunity.AllowedBranches); err != nil {
		log.Printf("Error unmarshaling allowed_branches: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(allowedGendersRaw, &opportunity.AllowedGenders); err != nil {
		log.Printf("Error unmarshaling allowed_genders: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Check eligibility
	criteria := utils.CheckEligibility(student, opportunity)
	isEligible := true
	for _, criterion := range criteria {
		if !criterion.Passed {
			isEligible = false
			break
		}
	}

	// If not eligible, return error
	if !isEligible {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"message":  "Student is not eligible for this opportunity.",
			"criteria": criteria,
		})
		return
	}

	// Create application record
	application := models.Application{
		StudentUSN:    student.USN,
		StudentName:   student.Name,
		OpportunityID: opportunity.ID,
		AppliedAt:     time.Now(),
	}
	
	tableName := "applications"
	if utils.CheckTableExists(db, tableName) {
		fmt.Printf("Table '%s' exists.\n", tableName)
	} else {
		fmt.Printf("Table '%s' does not exist. Creating table...\n", tableName)
		CreateApplicationsTable(db)
	}

	// Insert application into the database
	applyQuery := `INSERT INTO applications (student_usn, student_name, opportunity_id, applied_at) 
                   VALUES ($1, $2, $3, $4)`
	_, err = db.Exec(applyQuery, application.StudentUSN, application.StudentName, application.OpportunityID, application.AppliedAt)
	if err != nil {
		http.Error(w, "Error recording application", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Application submitted successfully.",
	})
}

func CreateApplicationsTable(db *sql.DB){
	query := `
		CREATE TABLE applications (
    	id SERIAL,
    	student_usn VARCHAR(20) PRIMARY KEY,
    	student_name VARCHAR(100) NOT NULL,
    	opportunity_id VARCHAR(20) NOT NULL,
    	applied_at TIMESTAMP NOT NULL,
		status VARCHAR(10) DEFAULT 'PROCESSING',
    	FOREIGN KEY (student_usn) REFERENCES students(usn),
    	FOREIGN KEY (opportunity_id) REFERENCES opportunities(id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	} else {
		log.Println("Table `applications` created or already exists.")
	}
}

func GetStudentApplicationsHandler(w http.ResponseWriter, r *http.Request,db *sql.DB){
	
	// Extract the student USN from query parameters
	studentUSN := r.URL.Query().Get("usn")
	if studentUSN == "" {
		http.Error(w, "Student USN is required", http.StatusBadRequest)
		return
	}

	// SQL query to fetch applications along with opportunity details
	query := `
		SELECT 
			applications.id, 
			applications.student_name, 
			applications.opportunity_id, 
			opportunities.title AS job_post, 
			opportunities.company, 
			applications.status
			FROM 
				applications 
			INNER JOIN 
				opportunities 
			ON 
				applications.opportunity_id = opportunities.id
			WHERE 
				applications.student_usn = $1
			ORDER BY 
				applications.applied_at DESC
		`
	// Execute the query
	rows, err := db.Query(query, studentUSN)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Struct to store application data
	type Application struct {
		ID             int    `json:"id"`
		StudentName    string `json:"student_name"`
		OpportunityID  string `json:"opportunity_id"`
		JobPost        string `json:"job_post"`
		Company        string `json:"company"`
		Status         string `json:"status"`
	}

	// Slice to store all applications
	var applications []Application

	for rows.Next() {
		var application Application
		if err := rows.Scan(&application.ID, &application.StudentName, &application.OpportunityID, &application.JobPost, &application.Company, &application.Status); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		applications = append(applications, application)
	}

	// Check if no applications were found
	if len(applications) == 0 {
		http.Error(w, "No applications found for the given student USN", http.StatusNotFound)
		return
	}

	// Return the applications as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applications)
}