package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"Github.com/Aryan-2511/Placement_NIE/models"
	"Github.com/Aryan-2511/Placement_NIE/utils"
)
type Coordinator struct {
	Name    string `json:"name"`
	Contact string `json:"contact"`
}
// GenerateOpportunityID creates unique IDs for opportunities: OP{YY}{001}
// YY represents last two digits of batch year
func GenerateOpportunityID(batch string, serial int) string {
	// Format serial number as a 3-digit string
	serialStr := fmt.Sprintf("%03d", serial)

	// Extract the last two digits of the batch (e.g., "2025" -> "25")
	batchCode := batch[len(batch)-2:]

	// Return formatted Opportunity-ID
	return fmt.Sprintf("OP%s%s", batchCode, serialStr)
}


func AddOpportunity(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	// log.Printf("Request Method: %s", r.Method)
	// log.Printf("Request Headers: %+v", r.Header)
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}
	var opportunity models.Opportunity
	if err := json.NewDecoder(r.Body).Decode(&opportunity); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if opportunity.Batch == "" {
		http.Error(w, "Batch is required", http.StatusBadRequest)
		return
	}

	allowedBranchesJSON, err := json.Marshal(opportunity.AllowedBranches)
	if err != nil {
		log.Printf("Error marshaling allowed_branches: %v", err)
		return
	}

	allowedGendersJSON, err := json.Marshal(opportunity.AllowedGenders)
	if err != nil {
		log.Printf("Error marshaling allowed_genders: %v", err)
		return
	}

	coordinatorsJSON, err := json.Marshal(opportunity.Coordinators)
	if err != nil {
		log.Printf("Error marshaling coordinators: %v", err)
		return
	}

	attachedDocumentsJSON, err := json.Marshal(opportunity.AttachedDocuments)
	if err != nil {
		log.Printf("Error marshaling attached_documents: %v", err)
		return
	}

	// Insert into database
	tableName := "opportunities"
	exists, err := utils.CheckTableExists(db, tableName)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if !exists {
		log.Printf("Table '%s' does not exist. Creating table...", tableName)
		CreateOpportunitiesTable(db) // Ensure this function is correctly implemented
	}

	var serial int
	query := `SELECT COUNT(*) + 1 FROM opportunities WHERE batch = $1`
	err = db.QueryRow(query, opportunity.Batch).Scan(&serial)
	if err != nil {
		log.Printf("Error fetching serial for Opportunity-ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	opportunityID := GenerateOpportunityID(opportunity.Batch, serial)

	query = `
	INSERT INTO opportunities (
		id, title, company, location, batch, ctc, ctc_info, cgpa, category, backlog, allowed_branches,
		allowed_genders, registration_date, coordinators, additional_info, form_link, job_description,
		attached_documents, opportunity_type, created_by, class_10_percentage_criteria, class_12_percentage_criteria, status
		) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9,
		$10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, 'ACTIVE'
		)`
	_, err = db.Exec(query,
		opportunityID,
		opportunity.Title,
		opportunity.Company,
		opportunity.Location,
		opportunity.Batch,
		opportunity.CTC,
		opportunity.CTCInfo,
		opportunity.CGPA,
		opportunity.Category,
		opportunity.Backlog,
		allowedBranchesJSON,
		allowedGendersJSON,
		opportunity.RegistrationDate,
		coordinatorsJSON,
		opportunity.AdditionalInfo,
		opportunity.FormLink,
		opportunity.JobDescription,
		attachedDocumentsJSON,
		opportunity.OpportunityType,
		claims["role"],
		opportunity.Class_10_Percentage_Criteria,
		opportunity.Class_12_Percentage_Criteria,
	)
	if err != nil {
		log.Printf("Error adding opportunity: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Opportunity added successfully"))
}

func CreateOpportunitiesTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS opportunities (
		id VARCHAR(10) PRIMARY KEY,
		title VARCHAR(255) NOT NULL,
		company VARCHAR(255) NOT NULL,
		location VARCHAR(255) NOT NULL,
		batch VARCHAR(10) NOT NULL,
		ctc NUMERIC(10, 2) NOT NULL,
		ctc_info TEXT,
		cgpa NUMERIC(3, 2) NOT NULL,
		category VARCHAR(50) NOT NULL,
		backlog INT NOT NULL,
		allowed_branches JSONB NOT NULL,
		allowed_genders JSONB NOT NULL,
		registration_date TIMESTAMP NOT NULL,
		coordinators JSONB NOT NULL,
		additional_info TEXT,
		form_link TEXT,
		job_description TEXT NOT NULL,
		attached_documents JSONB,
		opportunity_type VARCHAR(50) NOT NULL,
		created_by VARCHAR(255) NOT NULL,
		class_10_percentage_criteria NUMERIC(5, 2) DEFAULT 50,
		class_12_percentage_criteria NUMERIC(5, 2) DEFAULT 50,
		status VARCHAR(20) DEFAULT 'ACTIVE',
		completed VARCHAR(10) DEFAULT 'NO',
		created_at TIMESTAMP DEFAULT NOW() 
	);
	`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatalf("Error creating opportunities table: %v", err)
	} else {
		log.Println("Opportunities table ensured to exist.")
	}
}

func EditOpportunity(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Get Opportunity ID from query parameters
	opporID := r.URL.Query().Get("opportunity_id")
	if opporID == "" {
		http.Error(w, "Opportunity ID is required", http.StatusBadRequest)
		return
	}

	// Decode the request body into the Opportunity struct
	var opportunity models.Opportunity
	if err := json.NewDecoder(r.Body).Decode(&opportunity); err != nil {
		log.Printf("Error decoding request body: %v", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	allowedBranchesJSON, err := json.Marshal(opportunity.AllowedBranches)
	if err != nil {
		log.Printf("Error marshaling allowed_branches: %v", err)
		return
	}

	allowedGendersJSON, err := json.Marshal(opportunity.AllowedGenders)
	if err != nil {
		log.Printf("Error marshaling allowed_genders: %v", err)
		return
	}

	coordinatorsJSON, err := json.Marshal(opportunity.Coordinators)
	if err != nil {
		log.Printf("Error marshaling coordinators: %v", err)
		return
	}

	attachedDocumentsJSON, err := json.Marshal(opportunity.AttachedDocuments)
	if err != nil {
		log.Printf("Error marshaling attached_documents: %v", err)
		return
	}

	query := `
	UPDATE opportunities
	SET 
		title = $1,
		company = $2,
		location = $3,
		batch = $4,
		ctc = $5,
		ctc_info = $6,
		cgpa = $7,
		category = $8,
		backlog = $9,
		allowed_branches = $10,
		allowed_genders = $11,
		registration_date = $12,
		coordinators = $13,
		additional_info = $14,
		form_link = $15,
		job_description = $16,
		attached_documents = $17,
		opportunity_type = $18,
		created_by = $19,
		class_10_percentage_criteria = $20, 
		class_12_percentage_criteria = $21
	WHERE id = $22;
	`
	_, err = db.Exec(query,
		opportunity.Title,
		opportunity.Company,
		opportunity.Location,
		opportunity.Batch,
		opportunity.CTC,
		opportunity.CTCInfo,
		opportunity.CGPA,
		opportunity.Category,
		opportunity.Backlog,
		allowedBranchesJSON,
		allowedGendersJSON,
		opportunity.RegistrationDate,
		coordinatorsJSON,
		opportunity.AdditionalInfo,
		opportunity.FormLink,
		opportunity.JobDescription,
		attachedDocumentsJSON,
		opportunity.OpportunityType,
		claims["role"],
		opportunity.Class_10_Percentage_Criteria,
		opportunity.Class_12_Percentage_Criteria,
		opporID,
	)
	if err != nil {
		log.Printf("Error updating opportunity: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Opportunity updated successfully"))
}


func DeleteOpportunity(w http.ResponseWriter, r *http.Request,db *sql.DB){
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	opportunityId := r.URL.Query().Get("id")
	if opportunityId == ""{
		http.Error(w,"Opportunity ID is required", http.StatusBadRequest)
		return
	}

	query := `DELETE FROM opportunities WHERE id = $1`
	result,err := db.Exec(query,opportunityId)

	if err!=nil{
		log.Printf("Error deleting opportunity: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return	
	}
	rowsAffected,err := result.RowsAffected()
	if err!=nil{
		log.Printf("Error fetching rows affected: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0{
		http.Error(w, "Opportunity not found", http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Opportunity deleted successfully"))
}
func UpdateOpportunityCompletionStatus(db *sql.DB, opportunityID, completed string) error {
	if completed != "YES" && completed != "NO" {
		return fmt.Errorf("invalid value for completed field: Use 'YES' or 'NO'")
	}

	// Start database transaction
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start database transaction: %v", err)
	}
	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			log.Printf("Recovered from panic: %v", p)
		}
	}()

	// Update the opportunity's completed status
	updateOpportunityQuery := `
		UPDATE opportunities
		SET completed = $1
		WHERE id = $2
		RETURNING completed`
	var updatedStatus string
	err = tx.QueryRow(updateOpportunityQuery, completed, opportunityID).Scan(&updatedStatus)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update opportunity status: %v", err)
	}

	// If completed is set to "YES", close all associated applications
	if completed == "YES" {
		updateApplicationsQuery := `
			UPDATE applications
			SET status = 'CLOSED'
			WHERE opportunity_id = $1 AND status != 'CLOSED'`
		_, err := tx.Exec(updateApplicationsQuery, opportunityID)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update application status: %v", err)
		}
	}

	// Commit transaction
	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("failed to commit transaction: %v", err)
	}

	log.Printf("Opportunity completion status updated successfully: %s", updatedStatus)
	return nil
}
func UpdateOpportunityStatusHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
	
	type StatusUpdateResponse struct {
		UpdatedOpportunities int    `json:"updated_opportunities"`
		Message              string `json:"message"`
	}
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}


	query := `
		UPDATE opportunities
		SET status = 'CLOSED'
		WHERE registration_date <= $1
		AND status = 'ACTIVE';
		`
	now := time.Now()

	result, err := db.Exec(query, now)
	if err != nil {
		log.Printf("Error updating opportunity status: %v", err)
		http.Error(w, "Failed to update opportunity status", http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Printf("Error fetching rows affected: %v", err)
		http.Error(w, "Failed to fetch rows affected", http.StatusInternalServerError)
		return
	}

	response := StatusUpdateResponse{
		UpdatedOpportunities: int(rowsAffected),
		Message:              "Opportunity statuses updated successfully.",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}


func GetOpportunityDetailsHandler(w http.ResponseWriter, r *http.Request,db *sql.DB) {
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" && claims["role"]!="STUDENT"{
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}


	// Get the Opportunity ID from the request
	opportunityID := r.URL.Query().Get("id")
	if opportunityID == "" {
		http.Error(w, "Opportunity ID is required", http.StatusBadRequest)
		return
	}

	// Query the database for the opportunity details
	query := `
		SELECT id, title, company, location, batch, ctc, ctc_info, cgpa, category,
			backlog, allowed_branches, allowed_genders, registration_date, coordinators, 
			additional_info, form_link, job_description, attached_documents, 
			opportunity_type, created_by, class_10_percentage_criteria, 
			class_12_percentage_criteria, status
		FROM opportunities
		WHERE id = $1
	`
	var allowedBranchesRaw, allowedGendersRaw, coordinatorsRaw, attachedDocumentsRaw []byte
	var opportunity models.Opportunity
	err = db.QueryRow(query, opportunityID).Scan(
		&opportunity.ID, &opportunity.Title, &opportunity.Company, &opportunity.Location,
		&opportunity.Batch, &opportunity.CTC, &opportunity.CTCInfo, &opportunity.CGPA,
		&opportunity.Category, &opportunity.Backlog, &allowedBranchesRaw, &allowedGendersRaw,
		&opportunity.RegistrationDate, &coordinatorsRaw, &opportunity.AdditionalInfo,
		&opportunity.FormLink, &opportunity.JobDescription, &attachedDocumentsRaw,
		&opportunity.OpportunityType, &opportunity.Created_By,
		&opportunity.Class_10_Percentage_Criteria, &opportunity.Class_12_Percentage_Criteria,
		&opportunity.Status,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Opportunity not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
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

	if err := json.Unmarshal(coordinatorsRaw, &opportunity.Coordinators); err != nil {
		log.Printf("Error unmarshaling coordinators: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(attachedDocumentsRaw, &opportunity.AttachedDocuments); err != nil {
		log.Printf("Error unmarshaling attached_documents: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opportunity)
}

func GetOpportunitiesByBatchHandler(w http.ResponseWriter, r *http.Request, db *sql.DB) {
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
	if claims["role"] != "ADMIN" && claims["role"] != "PLACEMENT_COORDINATOR" && claims["role"] != "STUDENT" {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	// Automatically update the opportunity statuses before fetching
	queryUpdate := `
		UPDATE opportunities
		SET status = 'CLOSED'
		WHERE registration_date <= $1
		AND status = 'ACTIVE';`
	
	now := time.Now()
	_, err = db.Exec(queryUpdate, now)
	if err != nil {
		log.Printf("Error updating opportunity status: %v", err)
		http.Error(w, "Failed to update opportunity statuses", http.StatusInternalServerError)
		return
	}

	// Get the batch from the query parameter
	batch := r.URL.Query().Get("batch")
	if batch == "" {
		http.Error(w, "Batch is required", http.StatusBadRequest)
		return
	}

	// Query the database for opportunities for the specified batch
	queryFetch := `
		SELECT id, title, company, location, batch, ctc, category, registration_date, opportunity_type, status
		FROM opportunities
		WHERE batch = $1`
	
	rows, err := db.Query(queryFetch, batch)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var opportunities []struct {
		ID               string    `json:"id"`
		Title            string    `json:"title"`
		Company          string    `json:"company"`
		Location         string    `json:"location"`
		Batch            string    `json:"batch"`
		CTC              float64   `json:"ctc"`
		Category         string    `json:"category"`
		RegistrationDate time.Time `json:"registration_date"`
		Opportunity_type string    `json:"opportunity_type"`
		Status           string    `json:"status"`
	}

	for rows.Next() {
		var opportunity struct {
			ID               string    `json:"id"`
			Title            string    `json:"title"`
			Company          string    `json:"company"`
			Location         string    `json:"location"`
			Batch            string    `json:"batch"`
			CTC              float64   `json:"ctc"`
			Category         string    `json:"category"`
			RegistrationDate time.Time `json:"registration_date"`
			Opportunity_type string    `json:"opportunity_type"`
			Status           string    `json:"status"`
		}

		if err := rows.Scan(&opportunity.ID, &opportunity.Title, &opportunity.Company, &opportunity.Location, &opportunity.Batch, &opportunity.CTC, &opportunity.Category, &opportunity.RegistrationDate, &opportunity.Opportunity_type, &opportunity.Status); err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		opportunities = append(opportunities, opportunity)
	}

	// Return the list of opportunities as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(opportunities)
}
