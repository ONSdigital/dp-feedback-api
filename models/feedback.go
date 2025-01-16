package models

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-playground/validator/v10"
)

type Feedback struct {
	IsPageUseful      *bool  `json:"is_page_useful"         validate:"required"`
	IsGeneralFeedback *bool  `json:"is_general_feedback"    validate:"required"`
	OnsURL            string `json:"ons_url,omitempty"      validate:"omitempty,ons_url"`
	Feedback          string `json:"feedback,omitempty"`
	Name              string `json:"name,omitempty"`
	EmailAddress      string `json:"email_address,omitempty" validate:"omitempty,email"`
}

var cfg *config.Config

// getURLDomainValidator returns a validator func that checks that a field contains
// a valid URL with a hostname that ends with the provided domain
func getURLDomainValidator(onsDomain string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		urlField := fl.Field().String()
		return IsSiteDomainURL(urlField, onsDomain)
	}
}

// Validate checks that the Feedback struct complies with the validation tags
func (f *Feedback) Validate(cfg *config.Config) error {
	validate := validator.New()
	if err := validate.RegisterValidation("ons_url", getURLDomainValidator(cfg.OnsDomain)); err != nil {
		return fmt.Errorf("failed to register ons_url validator: %w", err)
	}
	return validate.Struct(f)
}

// Sanitize mutates all the strings in Feedback to prevent HTML, SQL and NoSQL injections,
// according to the provided config
func (f *Feedback) Sanitize(cfg *config.Sanitize) {
	f.OnsURL = Sanitize(cfg, f.OnsURL)
	f.Feedback = Sanitize(cfg, f.Feedback)
	f.Name = Sanitize(cfg, f.Name)
	f.EmailAddress = Sanitize(cfg, f.EmailAddress)
}

// IsSiteDomainURL is true when urlString is a URL and its host ends with `.`+siteDomain (when siteDomain is blank, or uses config.SiteDomain)
func IsSiteDomainURL(urlString, siteDomain string) bool {
	ctx := context.Background()
	if urlString == "" {
		return false
	}
	urlString = NormaliseURL(urlString)
	urlObject, err := url.ParseRequestURI(urlString)
	if err != nil {
		log.Error(ctx, "error parsing URL", err)
		return false
	}
	if siteDomain == "" {
		if cfg == nil {
			if cfg, err = config.Get(); err != nil {
				log.Error(ctx, "error getting config", err)
				return false
			}
		}
		siteDomain = cfg.OnsDomain
	}
	hostName := urlObject.Hostname()
	if hostName != siteDomain && !strings.HasSuffix(hostName, "."+siteDomain) {
		return false
	}
	return true
}

// NormaliseURL when a string is a URL without a scheme (e.g. `host.name/path`), add it (`https://`)
func NormaliseURL(urlString string) string {
	if strings.HasPrefix(urlString, "http") {
		return urlString
	}
	return "https://" + urlString
}
