package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)

// ApplyHandler processes student applications for opportunities
// Validates eligibility criteria and prevents duplicate applications
func ApplyHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	if claims["role"] != "STUDENT" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Parse request body
	var request struct {
		StudentUSN    string `json:"student_usn"`
		OpportunityID string `json:"opportunity_id"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Add validation for empty USN
	if request.StudentUSN == "" {
		log.Printf("Empty USN received in request")
		http.Error(w, "Student USN is required", http.StatusBadRequest)
		return
	}

	// Add debug log for the request payload
	log.Printf("Received application request - USN: %s, OpportunityID: %s", request.StudentUSN, request.OpportunityID)
	tableName := "applications"
	exists, err := utils.CheckTableExists(db, tableName)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Printf("Table '%s' does not exist. Creating table...", tableName)
		CreateApplicationsTable(db) // Ensure this function is correctly implemented
	}

	// Check if the student has already applied for the same opportunity
	var existingApplicationID int
	checkQuery := `SELECT id FROM applications WHERE student_usn = $1 AND opportunity_id = $2`
	err = db.QueryRow(checkQuery, request.StudentUSN, request.OpportunityID).Scan(&existingApplicationID)
	if err != nil && err != sql.ErrNoRows {
		log.Printf("Error checking existing application: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if existingApplicationID != 0 {
		http.Error(w, "You have already applied for this opportunity", http.StatusConflict)
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
	
	// Debug: Print the query and USN
	log.Printf("Executing query: %s with USN: %s", studentQuery, request.StudentUSN)
	
	// First, verify the exact data in the database
	var debugStudent struct {
	    USN string
	    Fields []string
	}
	debugRow := db.QueryRow("SELECT usn, array_agg(column_name::text) FROM students, information_schema.columns WHERE table_name='students' AND usn = $1 GROUP BY usn", request.StudentUSN)
	err = debugRow.Scan(&debugStudent.USN, &debugStudent.Fields)
	if err != nil {
	    log.Printf("Debug - Database content check failed: %v", err)
	} else {
	    log.Printf("Debug - Found student with USN: %s, Available fields: %v", debugStudent.USN, debugStudent.Fields)
	}

	// Original query
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
	    log.Printf("Detailed error when fetching student: %+v", err)
	    // Check column names in the database
	    var columnNames []string
	    colQuery := `
	        SELECT column_name 
	        FROM information_schema.columns 
	        WHERE table_name = 'students' 
	        ORDER BY ordinal_position;
	    `
	    rows, _ := db.Query(colQuery)
	    defer rows.Close()
	    for rows.Next() {
	        var colName string
	        rows.Scan(&colName)
	        columnNames = append(columnNames, colName)
	    }
	    log.Printf("Available columns in students table: %v", columnNames)
	    
	    http.Error(w, "Error fetching student details", http.StatusInternalServerError)
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

// CreateApplicationsTable initializes applications table with foreign keys
// Links students to opportunities with unique constraint
func CreateApplicationsTable(db *sql.DB) {
	query := `
		CREATE TABLE IF NOT EXISTS applications (
    	id SERIAL PRIMARY KEY,
    	student_usn VARCHAR(20) NOT NULL,
    	student_name VARCHAR(100) NOT NULL,
    	opportunity_id VARCHAR(20) NOT NULL,
    	applied_at TIMESTAMP NOT NULL,
    	status VARCHAR(10) DEFAULT 'IN-PROCESS',
    	FOREIGN KEY (student_usn) REFERENCES students(usn),
    	FOREIGN KEY (opportunity_id) REFERENCES opportunities(id),
    	CONSTRAINT unique_application UNIQUE (student_usn, opportunity_id)
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Failed to create table: %v", err)
	} else {
		log.Println("Table `applications` created or already exists.")
	}
}

// GetStudentApplicationsHandler retrieves all applications for a specific student
// Includes job post and company details from opportunities table
func GetStudentApplicationsHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
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
	if claims["role"] != "STUDENT" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
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

// GetApplicationsByBatch fetches all applications for students in a specific batch
// Returns extended application info including company and job details
func GetApplicationsByBatch(w http.ResponseWriter, r *http.Request, db *sql.DB) {
    if r.Method != http.MethodGet {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }

    // Authentication check
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

    claims, err := utils.ValidateToken(tokenString)
    if err != nil {
        log.Print(err)
        http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
        return
    }
    if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
        http.Error(w, "Unauthorized access", http.StatusForbidden)
        return
    }

    batch := r.URL.Query().Get("batch")
    if batch == "" {
        http.Error(w, "Batch parameter is required", http.StatusBadRequest)
        return
    }

    query := `
        SELECT 
            a.id,
            a.student_usn,
            a.student_name,
            a.opportunity_id,
            a.applied_at,
            a.status,
            o.company,
            o.title
        FROM 
            applications a
        INNER JOIN 
            students s ON a.student_usn = s.usn
        INNER JOIN 
            opportunities o ON a.opportunity_id = o.id
        WHERE 
            s.batch = $1
        ORDER BY 
            a.applied_at DESC`

    rows, err := db.Query(query, batch)
    if err != nil {
        log.Printf("Error fetching applications: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }
    defer rows.Close()

    type ExtendedApplication struct {
        models.Application
        Company string `json:"company"`
        JobTitle string `json:"job_title"`
    }

    var applications []ExtendedApplication
    for rows.Next() {
        var app ExtendedApplication
        err := rows.Scan(
            &app.ID,
            &app.StudentUSN,
            &app.StudentName,
            &app.OpportunityID,
            &app.AppliedAt,
            &app.Status,
            &app.Company,
            &app.JobTitle,
        )
        if err != nil {
            log.Printf("Error scanning application: %v", err)
            continue
        }
        applications = append(applications, app)
    }

    w.Header().Set("Content-Type", "application/json")
    if len(applications) == 0 {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "message": "No applications found for this batch",
            "applications": []ExtendedApplication{},
        })
        return
    }

    json.NewEncoder(w).Encode(map[string]interface{}{
        "batch": batch,
        "total_applications": len(applications),
        "applications": applications,
    })
}