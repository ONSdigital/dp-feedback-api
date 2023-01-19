package models

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"

	"github.com/ONSdigital/dp-feedback-api/config"
)

type Feedback struct {
	IsPageUseful      *bool  `json:"is_page_useful"`
	IsGeneralFeedback *bool  `json:"is_general_feedback"`
	OnsURL            string `json:"ons_url,omitempty"`
	Feedback          string `json:"feedback,omitempty"`
	Name              string `json:"name,omitempty"`
	EmailAddress      string `json:"email_address,omitempty"`
}

func (f *Feedback) Validate(cfg *config.Config) error {
	if f.IsPageUseful == nil {
		return errors.New("is_page_useful is compulsory")
	}

	if f.IsGeneralFeedback == nil {
		return errors.New("is_general_feedback is compulsory")
	}

	if f.EmailAddress != "" {
		if _, err := mail.ParseAddress(f.EmailAddress); err != nil {
			return fmt.Errorf("invalid email address: %w", err)
		}
	}

	if f.OnsURL != "" {
		onsURL, err := url.Parse(f.OnsURL)
		if err != nil {
			return fmt.Errorf("invalid ons url: %w", err)
		}

		if onsURL.Hostname() != cfg.OnsDomain {
			return fmt.Errorf("unexpected ons domain name: %s", onsURL.Hostname())
		}
	}

	return nil
}

func (f *Feedback) Sanitize(cfg *config.Sanitize) {
	f.OnsURL = Sanitize(cfg, f.OnsURL)
	f.Feedback = Sanitize(cfg, f.Feedback)
	f.Name = Sanitize(cfg, f.Name)
	f.EmailAddress = Sanitize(cfg, f.EmailAddress)
}
