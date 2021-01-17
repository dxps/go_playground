package main

import (
	"log"
	"os"

	"github.com/devisions/go-playground/go-send-email-options/01-net-smtp/pkg/email"
	"github.com/devisions/go-playground/go-send-email-options/01-net-smtp/pkg/smtp"
	"github.com/joho/godotenv"
)

type Config struct {
	smtpUsername string
	smtpPassword string
	from         string
	to           string
}

func main() {

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalln(">>> Failed to load the config. Reason:", err)
	}

	smtpClient := smtp.NewSMTPClientGMailTLS(cfg.smtpUsername, cfg.smtpPassword)

	msgBody := `Hello,<br/><br/>

	This is a test message.<br/><br/>

	Good day!
	`

	msg := email.NewMessage(cfg.from, cfg.to, "Test email sender", msgBody, true)

	if err := smtpClient.SendAndDone(*msg); err != nil {
		log.Println("[err] Failed to send the email. Reason:", err)
	}
}

func loadConfig() (*Config, error) {
	if err := godotenv.Load(".env"); err != nil {
		return nil, err
	}
	cfg := Config{
		smtpUsername: os.Getenv("SMTP_USERNAME"),
		smtpPassword: os.Getenv("SMTP_PASSWORD"),
		from:         os.Getenv("FROM"),
		to:           os.Getenv("TO"),
	}
	return &cfg, nil
}
