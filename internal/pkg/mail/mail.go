package mail

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"net/smtp"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
}

// LoadConfig loads SMTP configuration from environment variables
func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		logrus.Info("⚠️  Warning: .env file not found, using system env")
	}

	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))
	if port == 0 {
		port = 587 // default
	}

	return &Config{
		Host:     os.Getenv("SMTP_HOST"),
		Port:     port,
		Username: os.Getenv("SMTP_USERNAME"),
		Password: os.Getenv("SMTP_PASSWORD"),
		From:     os.Getenv("SMTP_FROM"),
	}
}

func SendWelcomeEmail(toEmail string) error {

	config := LoadConfig()

	if config.Host == "" || config.Username == "" || config.Password == "" {
		return fmt.Errorf("missing required SMTP configuration")
	}

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	from := config.From
	if from == "" {
		from = "admin@tokosukses.com" // default fallback
	}

	msg := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: Welcome to Our Store!\r\n" +
		"\r\n" +
		"Thank you for joining. Let's start shopping!\r\n")

	addr := config.Host + ":" + strconv.Itoa(config.Port)
	err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	logrus.Infof("✅ Welcome email sent successfully to %s", toEmail)
	return nil
}

func SendInvoiceEmail(toEmail string) error {

	config := LoadConfig()

	if config.Host == "" || config.Username == "" || config.Password == "" {
		return fmt.Errorf("missing required SMTP configuration")
	}

	auth := smtp.PlainAuth("", config.Username, config.Password, config.Host)

	from := config.From
	if from == "" {
		from = "admin@tokosukses.com"
	}

	msg := []byte("To: " + toEmail + "\r\n" +
		"From: " + from + "\r\n" +
		"Subject: Payment Received!\r\n" +
		"\r\n" +
		"Thank you! We have received your payment. Your order is now being processed.\r\n")

	addr := config.Host + ":" + strconv.Itoa(config.Port)
	err := smtp.SendMail(addr, auth, from, []string{toEmail}, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	logrus.Infof("✅ Invoice email sent successfully to %s", toEmail)
	return nil
}
