package utils

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

func SendBulkEmail(recipients []string, subject, body string) error {
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

	d := gomail.NewDialer("smtp.gmail.com", 587, email, password)
	var wg sync.WaitGroup
	errorChannel := make(chan error, len(recipients))

	for _, recipient := range recipients {
		wg.Add(1)
		go func(recipient string) {
			defer wg.Done()
			m := gomail.NewMessage()
			m.SetHeader("From", email)
			m.SetHeader("To", recipient)
			m.SetHeader("Subject", subject)
			m.SetBody("text/html", body)

			if err := d.DialAndSend(m); err != nil {
				log.Printf("Error sending email to %s: %v\n", recipient, err)
				errorChannel <- err
			}
		}(recipient)
	}

	wg.Wait()
	close(errorChannel)

	if len(errorChannel) > 0 {
		return fmt.Errorf("some emails failed to send")
	}
	return nil
}
