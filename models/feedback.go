package models

import (
	"fmt"
	"net/mail"
	"net/url"
	"strings"

	"github.com/ONSdigital/dp-feedback-api/config"
)

type Feedback struct {
	IsPageUseful      bool   `json:"is_page_useful"`
	IsGeneralFeedback bool   `json:"is_general_feedback"`
	OnsURL            string `json:"ons_url,omitempty"`
	Feedback          string `json:"feedback,omitempty"`
	Name              string `json:"name,omitempty"`
	EmailAddress      string `json:"email_address,omitempty"`
}

func (f *Feedback) Validate(cfg *config.Config) error {
	if _, err := mail.ParseAddress(f.EmailAddress); err != nil {
		return fmt.Errorf("invalid email address: %w", err)
	}

	onsURL, err := url.Parse(f.OnsURL)
	if err != nil {
		return fmt.Errorf("invalid ons url: %w", err)
	}

	if onsURL.Hostname() != cfg.OnsDomain {
		return fmt.Errorf("unexpected ons domain name: %s", onsURL.Hostname())
	}

	if err := sanitize(f.OnsURL, f.Feedback, f.Name, f.EmailAddress); err != nil {
		return fmt.Errorf("sanitization error: %w", err)
	}
	return nil
}

func sanitize(strs ...string) error {
	// TODO check for html injection, sql injection, nosql injection
	for _, toSanitize := range strs {
		if strings.ContainsAny(toSanitize, "&^%") {
			return fmt.Errorf("%s contains invalid characters", toSanitize)
		}
	}
	return nil
}
