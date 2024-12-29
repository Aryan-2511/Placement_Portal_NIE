package utils
import (
	"database/sql"
	"log"
)

func CheckTableExists(db *sql.DB, tableName string) bool {
	query := `
		SELECT EXISTS (
			SELECT 1
			FROM information_schema.tables 
			WHERE table_schema = 'public'
			AND table_catalog = 'Placement_Portal'
			AND table_name = $1
		);
	`

	var exists bool
	err := db.QueryRow(query, tableName).Scan(&exists)
	if err != nil {
		log.Fatalf("Error checking table existence: %v", err)
	}

	return exists
}