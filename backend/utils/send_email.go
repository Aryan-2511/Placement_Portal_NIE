package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	err := godotenv.Load()
    if err != nil {
        log.Fatalf("Error loading .env file: %v", err)
    }
	email := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWORD")

	if email == "" || password == "" {
		log.Println("Error: Email credentials are not set in environment variables")
		return fmt.Errorf("email credentials are missing")
	}
	m := gomail.NewMessage()
	m.SetHeader("From", email) // Replace with your email
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	d := gomail.NewDialer("smtp.gmail.com", 587, email, password	) 

	if err := d.DialAndSend(m); err != nil {
		log.Printf("Error sending email: %v\n", err)
		return err
	}
	return nil
}
