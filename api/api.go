package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/middleware"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
)

// API provides a struct to wrap the api around
type API struct {
	Cfg         *config.Config
	Router      chi.Router
	EmailSender EmailSender
}

// Setup function sets up the api and returns an api
func Setup(ctx context.Context, cfg *config.Config, r chi.Router, e EmailSender) *API {
	api := &API{
		Cfg:         cfg,
		Router:      r,
		EmailSender: e,
	}

	api.mountEndpoints(ctx)

	return api
}

// mountEndpoints creates a a new chi Router with the auth middleware and required endpoints,
// and then mounts it to the existing router, in order to prevent existing endpoints (i.e. /health) to go through auth.
func (api *API) mountEndpoints(ctx context.Context) {
	r := chi.NewRouter()
	middleware.Setup(r, api.Cfg.ServiceAuthToken)

	r.Post("/feedback", api.PostFeedback)

	api.Router.Mount("/", r)
}

// unmarshal is an aux function to read the provided ReadCloser and unmarshal it to the provided model struct
func Unmarshal(body io.ReadCloser, v interface{}) error {
	b, err := io.ReadAll(body)
	if err != nil {
		return fmt.Errorf("failed to read req body: %w", err)
	}

	if err = json.Unmarshal(b, v); err != nil {
		return fmt.Errorf("failed to unmarshal req body into a model: %s", err)
	}
	return nil
}

// WriteJSON responds generates
func (api *API) WriteJSON(ctx context.Context, w http.ResponseWriter, status int, resp interface{}) error {
	b, err := json.Marshal(resp)
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err = w.Write(b); err != nil {
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}

func (api *API) handleError(ctx context.Context, w http.ResponseWriter, err error, status int) {
	log.Error(ctx, "request failed", err)
	http.Error(w, err.Error(), status)
}
