package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	// "context"

	// "github.com/aws/aws-sdk-go-v2/config"
	// "github.com/aws/aws-sdk-go-v2/feature/rds/auth"
	_ "github.com/lib/pq"
)
var DB *sql.DB
func InitDB() *sql.DB{
	
	// // Load environment variables
	// dbUser := os.Getenv("DB_USER")
	// dbPassword := os.Getenv("DB_PASSWORD")
	// dbName := os.Getenv("DB_NAME")
	// dbSSLMode := os.Getenv("DB_SSL_MODE")

	// if dbUser == "" || dbPassword == "" || dbName == "" || dbSSLMode == "" {
	// 	log.Fatalf("Database credentials are not set in environment variables")
	// }

	// // Construct the connection string
	// connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dbUser, dbPassword, dbName, dbSSLMode)

	// // Open a connection to the database
	// DB, err = sql.Open("postgres", connStr)
	// if err != nil {
	// 	log.Fatalf("Failed to connect to database: %v", err)
	// }

	// // Verify the connection
	// if err = DB.Ping(); err != nil {
	// 	log.Fatalf("Failed to ping the database: %v", err)
	// }

	// fmt.Println("Successfully connected to the database")
	// return DB
	var err error

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Fetch database connection details from environment variables
	dbURL := os.Getenv("RENDER_DB")
	if dbURL == "" {
		log.Fatal("DATABASE_URL is not set in environment variables")
	}

	// Open a connection to the database
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Verify the connection
	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the Render PostgreSQL database!")
	return DB
}
