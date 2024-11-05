package models

import (
	"fmt"
	"net/url"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/go-playground/validator/v10"
)

type Feedback struct {
	IsPageUseful      *bool  `json:"is_page_useful"         validate:"required"`
	IsGeneralFeedback *bool  `json:"is_general_feedback"    validate:"required"`
	OnsURL            string `json:"ons_url,omitempty"      validate:"omitempty,ons_url"`
	Feedback          string `json:"feedback"				validate:"required"`
	Name              string `json:"name,omitempty"`
	EmailAddress      string `json:"email_address,omitempty" validate:"omitempty,email"`
}

// getURLDomainValidator returns a validator func that checks that a field contains
// a valid URL with a hostname corresponding to the provided domain
func getURLDomainValidator(domain string) validator.Func {
	return func(fl validator.FieldLevel) bool {
		val := fl.Field().String()
		onsURL, err := url.Parse(val)
		if err != nil {
			return false
		}
		return onsURL.Hostname() == domain
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
