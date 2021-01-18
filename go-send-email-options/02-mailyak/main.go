package main

import (
	"log"
	"net/smtp"
	"os"

	"github.com/domodwyer/mailyak/v3"
	"github.com/joho/godotenv"
)

type Config struct {
	smtpUsername string
	smtpPassword string
	from         string // email address of the sender
	fromName     string // name of the sender
	to           string // email address of the recipient
}

func main() {

	cfg, err := loadConfig()
	if err != nil {
		log.Fatalln(">>> Failed to load the config. Reason:", err)
	}

	msgSubject := "Test email sender (from)"
	msgBody := `Hello,<br/><br/>

	This is a test message.<br/><br/>

	Good day!
	`

	mail, err := mailyak.NewWithTLS(
		"smtp.gmail.com:465",
		smtp.PlainAuth("", cfg.smtpUsername, cfg.smtpPassword, "smtp.gmail.com"),
		nil,
	)
	if err != nil {
		log.Fatalln(">>> Failed to init. Reason:", err)
	}

	mail.From(cfg.from)
	mail.FromName(cfg.fromName)
	mail.To(cfg.to)
	mail.Subject(msgSubject)
	_, err = mail.HTML().Write([]byte(msgBody))
	if err != nil {
		log.Fatalln(">>> Failed to setup content. Reason:", err)
	}
	if err := mail.Send(); err != nil {
		log.Fatalln(">>> Failed to send email. Reason:", err)
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
		fromName:     os.Getenv("FROM_NAME"),
		to:           os.Getenv("TO"),
	}
	return &cfg, nil
}
