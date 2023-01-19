package middleware

import (
	"errors"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/headers"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
)

// Setup adds the middleware chain to the provided router
func Setup(r chi.Router, serviceAuthToken string) chi.Router {
	checkServiceIdentity := serviceAuthHandler(serviceAuthToken)
	r.Use(checkServiceIdentity)
	return r
}

// serviceAuthHandler returns a handler to perform service authentication.
func serviceAuthHandler(serviceAuthToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()

			token, err := headers.GetServiceAuthToken(r)
			if err != nil {
				log.Error(ctx, "auth failed", err)
				http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
				return
			}

			// TODO At the moment it only checks the token matches the expected config value,
			// but once the new service authentication is implemented, it should be used here.
			if token != serviceAuthToken {
				log.Error(ctx, "auth failed", errors.New("wrong service auth token provided"))
				http.Error(w, "401 Unauthorized", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
