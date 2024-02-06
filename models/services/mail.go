package services

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-mail/mail/v2"
)

const (
	DefaultEmail string = "support@lenslocked.com"
)

var (
	ErrFailedToSendEmail              = errors.New("failed to send e-mail")
	ErrFailedToSendResetPasswordEmail = errors.New("failed to send reset password e-mail")
)

type Email struct {
	From      string
	To        string
	Subject   string
	PlainText string
	HTML      string
}

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

func (es *EmailService) Send(email Email) error {
	msg := mail.NewMessage()
	msg.SetHeader("From", es.getEmailFromOrDefault(email.From))
	msg.SetHeader("To", email.To)
	msg.SetHeader("Subject", email.Subject)
	pt := strings.TrimSpace(email.PlainText)
	html := strings.TrimSpace(email.HTML)
	if pt != "" && html != "" {
		msg.SetBody("text/plain", email.PlainText)
		msg.AddAlternative("text/html", email.HTML)
	} else if pt != "" {
		msg.SetBody("text/plain", email.PlainText)
	} else if html != "" {
		msg.SetBody("text/plain", email.PlainText)
	}

	err := es.dialer.DialAndSend(msg)
	if err != nil {
		return errors.Join(ErrFailedToSendEmail, err)
	}
	return nil
}

func (es *EmailService) getEmailFromOrDefault(from string) string {
	if from := strings.TrimSpace(from); from != "" {
		return from
	} else if de := strings.TrimSpace(es.DefaultEmail); de != "" {
		return de
	}
	return DefaultEmail
}

func (es *EmailService) ForgotPassword(to, resetURL string) error {
	err := es.Send(Email{
		From:      "",
		To:        to,
		Subject:   "Reset your password",
		PlainText: fmt.Sprintf("To reset your password, please visit the following link: %s", resetURL),
		HTML: fmt.Sprintf(`<p><To reset your password, please visit the following link: <a href="%s">%[1]s</a></p>`,
			resetURL),
	})
	if err != nil {
		return errors.Join(ErrFailedToSendResetPasswordEmail, err)
	}

	return nil
}
