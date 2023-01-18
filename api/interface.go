package api

//go:generate moq -out mock/email.go -pkg mock . EmailSender

// EmailSender defines the required methods from the email sender package
type EmailSender interface {
	Send(from string, to []string, msg []byte) error
}
