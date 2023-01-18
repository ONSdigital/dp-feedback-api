package middleware

import (
	"github.com/ONSdigital/dp-api-clients-go/v2/identity"
	"github.com/ONSdigital/dp-net/v2/handlers"
	"github.com/go-chi/chi/v5"
)

// UseAuth adds a middleware to check service identity only
func UseAuth(r chi.Router, idc *identity.Client) chi.Router {
	checkServiceIdentity := handlers.IdentityService(idc)
	r.Use(checkServiceIdentity)
	return r
}
