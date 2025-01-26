package utils

import (
	"database/sql"
	"fmt"
)

func CheckTableExists(db *sql.DB, tableName string) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public'
			AND table_catalog = current_database() -- Uses current database dynamically
			AND table_name = $1
		);
	`

	var exists bool
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("error checking table existence: %v", err)
	}

	return exists, nil
}