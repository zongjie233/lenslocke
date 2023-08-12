package models

import "github.com/go-mail/mail/v2"

const (
	DefaultSender = "support@zdq.com"
)

type SMTPConfig struct {
	host     string
	port     int
	username string
	password string
}

func NewEmailService(config SMTPConfig) *EmailService {
	es := EmailService{
		dialer: mail.NewDialer(config.host, config.port, config.username, config.password),
	}
	return &es
}

type EmailService struct {
	DefaultSender string

	dialer *mail.Dialer
}
