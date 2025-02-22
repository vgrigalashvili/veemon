package mail

import (
	"context"
)

const (
	smtpAuthAddress   = "smtp.gmail.com"
	smtpServerAddress = "smtp.gmail.com:587"
)

// EmailSender defines an interface for sending emails.
type EmailSender interface {
	SendEmail(ctx context.Context, to []string, subject, body string) error
}

// type emailSender struct {
// 	name              string
// 	fromEmailAddress  string
// 	fromEmailPassword string
// }

// func NewEmailSender(ctx context.Context, name string, fromEmailAddress string, fromEmailPassword string) EmailSender {
// 	return &emailSender{
// 		name:              name,
// 		fromEmailAddress:  fromEmailAddress,
// 		fromEmailPassword: fromEmailPassword,
// 	}
// }

// func (sender *emailSender) SendEmail(
// 	ctx context.Context,
// 	to []string,
// 	subject string,
// 	body string,
// ) error {
// 	e := email.NewEmail()
// 	e.From = fmt.Sprintf("%s <%s>", sender.name, sender.fromEmailAddress)
// 	e.Subject = subject
// 	e.HTML = []byte(body)
// 	e.To = to

// 	smtpAuth := smtp.PlainAuth("", sender.fromEmailAddress, sender.fromEmailPassword, smtpAuthAddress)
// 	return e.Send(smtpServerAddress, smtpAuth)
// }
