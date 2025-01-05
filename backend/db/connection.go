package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"

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
	if err := godotenv.Load(); err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
	dbName  := os.Getenv("DB_NAME") 
	dbUser 	:= os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPortStr := os.Getenv("DB_PORT")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		log.Fatalf("Invalid integer for dbPort: %v", err)
	}


	// Data Source Name (DSN) with password authentication
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=require",
		dbHost, dbPort, dbUser, dbPassword, dbName,
	)

	// Connect to the database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}

	// Verify the connection
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
	return db
}
