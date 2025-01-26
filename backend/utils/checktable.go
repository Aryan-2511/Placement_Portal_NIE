package utils

import (
	"database/sql"
	"fmt"
	"log"
)

func CheckTableExists(db *sql.DB, tableName string) (bool, error) {
	query := fmt.Sprintf(`
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public'
			AND table_name = '%s'
		);
	`, tableName)

	var exists bool
	err := db.QueryRow(query).Scan(&exists)
	if err != nil {
		log.Printf("Error checking table existence: %v", err)
		return false, err
	}

	return exists, nil
}