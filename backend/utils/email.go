package utils

import (
	"database/sql"
	"log"
	"os"

	"github.com/go-mail/mail/v2"
)

func SendEmailsAsync(recipients []string, subject, message, notificationID string, db *sql.DB) {
	// Get email credentials from environment variables
	senderEmail := os.Getenv("EMAIL_ADDRESS")
	appPassword := os.Getenv("EMAIL_APP_PASSWORD")

	if senderEmail == "" || appPassword == "" {
		log.Printf("Email credentials not configured")
		return
	}

	dialer := mail.NewDialer("smtp.gmail.com", 587, senderEmail, appPassword)

	for _, recipient := range recipients {
		m := mail.NewMessage()
		m.SetHeader("From", senderEmail)
		m.SetHeader("To", recipient)
		m.SetHeader("Subject", subject)
		m.SetBody("text/html", message)

		if err := dialer.DialAndSend(m); err != nil {
			log.Printf("Failed to send email to %s: %v", recipient, err)
			continue
		}

		// Update sent status
		_, err := db.Exec(`
            UPDATE notification_recipients 
            SET is_sent = TRUE, sent_at = CURRENT_TIMESTAMP 
            WHERE notification_id = $1 AND recipient_email = $2`,
			notificationID, recipient)
		if err != nil {
			log.Printf("Failed to update sent status: %v", err)
		}
	}
}
