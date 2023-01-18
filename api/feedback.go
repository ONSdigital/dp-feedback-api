package api

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-feedback-api/models"
)

// PostFeedback returns a handler for POST /feedback

func (api *API) PostFeedback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feedback := &models.Feedback{}
	if err := unmarshal(r.Body, feedback); err != nil {
		api.handleError(ctx, w, err, http.StatusBadRequest)
		return
	}

	if err := feedback.Validate(api.Cfg); err != nil {
		api.handleError(ctx, w, err, http.StatusBadRequest)
		return
	}

	if err := api.EmailSender.Send(
		api.Cfg.Mail.From,
		[]string{api.Cfg.Mail.To},
		generateFeedbackMessage(feedback, api.Cfg.Mail.From, api.Cfg.Mail.To),
	); err != nil {
		api.handleError(ctx, w, fmt.Errorf("failed to send message: %w", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func generateFeedbackMessage(f *models.Feedback, from, to string) []byte {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("From: %s\n", from))
	b.WriteString(fmt.Sprintf("To: %s\n", to))
	b.WriteString(fmt.Sprintf("Subject: Feedback received\n\n"))

	b.WriteString(fmt.Sprintf("Is page useful: %t\n", *f.IsPageUseful))
	b.WriteString(fmt.Sprintf("Is general feedback: %t\n", *f.IsGeneralFeedback))
	b.WriteString(fmt.Sprintf("Page URL: %s\n", f.OnsURL))
	b.WriteString(fmt.Sprintf("Description: %s\n", f.Feedback))

	if len(f.Name) > 0 {
		b.WriteString(fmt.Sprintf("Name: %s\n", f.Name))
	}

	if len(f.EmailAddress) > 0 {
		b.WriteString(fmt.Sprintf("Email address: %s\n", f.EmailAddress))
	}

	return b.Bytes()
}
