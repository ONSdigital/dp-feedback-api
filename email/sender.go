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
	var auth smtp.Auth
	if cfg.Encrypted {
		auth = smtp.PlainAuth(
			"",
			cfg.User,
			cfg.Password,
			cfg.Host,
		)
	} else {
		auth = smtp.CRAMMD5Auth(cfg.User, cfg.Password)
	}
	mailAddr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	return &SMTPSender{
		Addr: mailAddr,
		Auth: auth,
	}
}
