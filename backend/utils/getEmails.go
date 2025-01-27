package utils

import (
	"database/sql"
)
func GetEmailsByBatch(db *sql.DB, batch string) ([]string, error) {
	query := `SELECT college_email FROM students WHERE batch = $1`
	rows, err := db.Query(query, batch)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}
func GetEmailsByOpportunity(db *sql.DB, opportunityID string) ([]string, error) {
	query := `
		SELECT s.college_email
		FROM students s
		INNER JOIN applications a
		ON s.usn = a.usn
		WHERE a.opportunity_id = $1
	`
	rows, err := db.Query(query, opportunityID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emails []string
	for rows.Next() {
		var email string
		if err := rows.Scan(&email); err != nil {
			return nil, err
		}
		emails = append(emails, email)
	}
	return emails, nil
}
