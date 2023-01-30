package api

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-feedback-api/models"
)

// PostFeedback is the handler for POST /feedback
// It unmarshals and validates the feedback data before sending to configured email account
func (api *API) PostFeedback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	feedback := &models.Feedback{}
	if err := Unmarshal(r.Body, feedback); err != nil {
		api.handleError(ctx, w, err, http.StatusBadRequest)
		return
	}

	if err := feedback.Validate(api.Cfg); err != nil {
		api.handleError(ctx, w, err, http.StatusBadRequest)
		return
	}

	feedback.Sanitize(api.Cfg.Sanitize)

	if err := api.EmailSender.Send(
		api.Cfg.FeedbackFrom,
		[]string{api.Cfg.FeedbackTo},
		GenerateFeedbackMessage(feedback, api.Cfg.FeedbackFrom, api.Cfg.FeedbackTo),
	); err != nil {
		api.handleError(ctx, w, fmt.Errorf("failed to send message: %w", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func GenerateFeedbackMessage(f *models.Feedback, from, to string) []byte {
	var b bytes.Buffer

	b.WriteString(fmt.Sprintf("From: %s\n", from))
	b.WriteString(fmt.Sprintf("To: %s\n", to))
	b.WriteString(fmt.Sprintf("Subject: Feedback received\n\n"))

	b.WriteString(fmt.Sprintf("Is page useful: %t\n", *f.IsPageUseful))
	b.WriteString(fmt.Sprintf("Is general feedback: %t\n", *f.IsGeneralFeedback))

	if len(f.OnsURL) > 0 {
		b.WriteString(fmt.Sprintf("Page URL: %s\n", f.OnsURL))
	}

	if len(f.Feedback) > 0 {
		b.WriteString(fmt.Sprintf("Description: %s\n", f.Feedback))
	}

	if len(f.Name) > 0 {
		b.WriteString(fmt.Sprintf("Name: %s\n", f.Name))
	}

	if len(f.EmailAddress) > 0 {
		b.WriteString(fmt.Sprintf("Email address: %s\n", f.EmailAddress))
	}

	return b.Bytes()
}
