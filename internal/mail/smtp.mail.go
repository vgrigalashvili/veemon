package mail

import (
	"context"
	"fmt"
	"net/smtp"
)

type SMTPMailer struct {
	host     string
	port     string
	username string
	password string
	from     string
}

// NewSMTPMailer creates a new SMTPMailer instance.
func NewSMTPMailer(host, port, username, password, from string) *SMTPMailer {
	return &SMTPMailer{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// SendEmail sends an email using the SMTP server.
func (s *SMTPMailer) SendEmail(ctx context.Context, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	msg := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\n\r\n%s", s.from, to, subject, body))

	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		if err := smtp.SendMail(addr, auth, s.from, to, msg); err != nil {
			return fmt.Errorf("failed to send email: %w", err)
		}
	}

	return nil
}
