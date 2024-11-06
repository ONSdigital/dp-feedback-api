package email

import (
	"fmt"
	"net/smtp"

	"github.com/ONSdigital/dp-feedback-api/config"
)

type SMTPSender struct {
	Addr string
	Auth smtp.Auth
}

type unencryptedAuth struct {
	smtp.Auth
}

func (s SMTPSender) Send(from string, to []string, msg []byte) error {
	return smtp.SendMail(
		s.Addr,
		s.Auth,
		from,
		to,
		msg)
}

// NewSMTPSender returns a new SMTPSender according to the provided mail configuration
func NewSMTPSender(cfg *config.Mail) *SMTPSender {
	auth := smtp.PlainAuth(
		"",
		cfg.User,
		cfg.Password,
		cfg.Host,
	)
	if cfg.Host == "localhost" {
		auth = unencryptedAuth{auth}
	}
	mailAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &SMTPSender{
		Addr: mailAddr,
		Auth: auth,
	}
}
