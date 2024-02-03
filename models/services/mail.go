package services

import (
	"github.com/go-mail/mail/v2"
)

const (
	DefaultEmail string = "support@lenslocked.com"
)

type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
}

type EmailService struct {
	DefaultEmail string
	dialer       *mail.Dialer
}

func NewEmailService(config SMTPConfig) *EmailService {
	return &EmailService{
		dialer: mail.NewDialer(config.Host, config.Port, config.Username, config.Password),
	}
}
