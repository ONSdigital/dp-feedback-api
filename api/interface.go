package api

type EmailSender interface {
	Send(from string, to []string, msg []byte) error
}
