package mail

import (
	"fmt"
	"net/smtp"

	"github.com/jordan-wright/email"
)

type EmailSender interface {
	SendEmail(
		subject string,
		content string,
		to []string,
		cc []string,
		bcc []string,
		attachedFiles []string,
	) error
}

type GmailSender struct {
	name      string
	fromEmail string
	fromPsswd string
}

// SendEmail implements EmailSender.
func (s *GmailSender) SendEmail(subject string, content string, to []string, cc []string, bcc []string, attachedFiles []string) error {
	e := email.NewEmail()

	e.From = fmt.Sprintf("%s <%s>", s.name, s.fromEmail)
	e.Subject = subject
	e.HTML = []byte(content)
	e.To = to
	e.Cc = cc
	e.Bcc = bcc

	for _, f := range attachedFiles {
		_, err := e.AttachFile(f)
		if err != nil {
			return fmt.Errorf("failed to attach file %s: %w", f, err)
		}
	}

	smtpAuth := smtp.PlainAuth("", s.fromEmail, s.fromPsswd, "smtp.gmail.com")
	return e.Send("smtp.gmail.com:587", smtpAuth)
}

func NewEmailSender(name string, fromEmail string, fromPsswd string) EmailSender {
	return &GmailSender{
		name:      name,
		fromEmail: fromEmail,
		fromPsswd: fromPsswd,
	}
}
