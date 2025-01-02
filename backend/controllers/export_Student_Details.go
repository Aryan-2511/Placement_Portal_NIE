package controllers

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"Github.com/Aryan-2511/Placement_NIE/db"
)

func ExportCustomStudentDetailsToCSV(w http.ResponseWriter, r *http.Request) {
	userRole := r.Header.Get("Role")
	if userRole != "ADMIN" && userRole != "PLACEMENT_COORDINATOR" {
		http.Error(w, "Unauthorized: Only admins or placement coordinators can add opportunities", http.StatusUnauthorized)
		return
	}
	var request struct {
		OpportunityID string   `json:"opportunity_id"`
		Fields        []string `json:"fields"` // Fields to export
	}

	// Decode request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Validate required parameters
	if request.OpportunityID == "" {
		http.Error(w, "Opportunity ID is required", http.StatusBadRequest)
		return
	}
	if len(request.Fields) == 0 {
		http.Error(w, "At least one field must be selected for export", http.StatusBadRequest)
		return
	}

	// List all allowed fields from the `students` table
	allowedFields := map[string]string{
		"usn":                  "students.usn",
		"name":                 "students.name",
		"dob":                  "students.dob",
		"college_email":        "students.college_email",
		"personal_email":       "students.personal_email",
		"branch":               "students.branch",
		"batch":                "students.batch",
		"address":              "students.address",
		"contact":              "students.contact",
		"gender":               "students.gender",
		"category":             "students.category",
		"aadhar":               "students.aadhar",
		"pan":                  "students.pan",
		"class_10_percentage":  "students.class_10_percentage",
		"class_10_year":        "students.class_10_year",
		"class_10_board":       "students.class_10_board",
		"class_12_percentage":  "students.class_12_percentage",
		"class_12_year":        "students.class_12_year",
		"class_12_board":       "students.class_12_board",
		"current_cgpa":         "students.current_cgpa",
		"backlogs":             "students.backlogs",
	}

	// Validate selected fields and build SQL query dynamically
	var selectedColumns []string
	for _, field := range request.Fields {
		if col, ok := allowedFields[field]; ok {
			selectedColumns = append(selectedColumns, col+" AS "+field)
		} else {
			http.Error(w, fmt.Sprintf("Invalid field: %s", field), http.StatusBadRequest)
			return
		}
	}

	query := fmt.Sprintf(`
		SELECT %s
		FROM applications
		INNER JOIN students ON applications.student_usn = students.usn
		WHERE applications.opportunity_id = $1
		ORDER BY students.usn
	`, strings.Join(selectedColumns, ", "))
	log.Print(query)
	// Execute the query
	rows, err := db.DB.Query(query, request.OpportunityID)
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Get column names from the selected fields
	columnNames := request.Fields

	// Set headers for file download
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=custom_applications_opportunity_%s.csv", request.OpportunityID))
	w.Header().Set("Content-Type", "text/csv")

	writer := csv.NewWriter(w)
	defer writer.Flush()

	// Write headers to the CSV file
	if err := writer.Write(columnNames); err != nil {
		log.Printf("Error writing CSV headers: %v", err)
		http.Error(w, "Error generating file", http.StatusInternalServerError)
		return
	}

	// Write rows to the CSV file
	for rows.Next() {
		values := make([]interface{}, len(columnNames))
		valuePtrs := make([]interface{}, len(columnNames))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			log.Printf("Row scan error: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		record := make([]string, len(columnNames))
	for i, val := range values {
    	if val != nil {
        	// Handle DECIMAL(4,2) fields
        	if decimalValue, ok := val.([]byte); ok {
            	record[i] = string(decimalValue) // Convert byte array to string
        	} else {
            	record[i] = fmt.Sprintf("%v", val)
        	}
    	} else {
        	record[i] = ""
    	}
	}

		if err := writer.Write(record); err != nil {
			log.Printf("Error writing CSV record: %v", err)
			http.Error(w, "Error generating file", http.StatusInternalServerError)
			return
		}
	}
}
