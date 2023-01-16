package api

import (
	"net/http"

	"github.com/ONSdigital/dp-feedback-api/models"
	"github.com/ONSdigital/log.go/v2/log"
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

	log.Info(ctx, "OK")
}
