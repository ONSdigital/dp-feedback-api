package models_test

import (
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	onsHost           = "testhost"
	cfg               = &config.Config{OnsDomain: onsHost}
	pageIsUseful      = true
	isGeneralFeedback = true
)

func validFeedbackModel() *models.Feedback {
	return &models.Feedback{
		IsPageUseful:      &pageIsUseful,
		IsGeneralFeedback: &isGeneralFeedback,
		OnsURL:            fmt.Sprintf("https://%s:1234/sub/path", onsHost),
		Feedback:          "very nice and useful website!",
		Name:              "Mr Feedback reporter",
		EmailAddress:      "feedback@reporter.com",
	}
}

func TestValidate(t *testing.T) {
	Convey("Given a valid fully populated Feedback model", t, func() {
		f := validFeedbackModel()

		Convey("Then validation is successful", func() {
			So(f.Validate(cfg), ShouldBeNil)
		})
	})

	Convey("Given a valid Feedback model containing only the mandatory fields", t, func() {
		f := &models.Feedback{
			IsPageUseful:      &pageIsUseful,
			IsGeneralFeedback: &isGeneralFeedback,
		}

		Convey("Then validation is successful", func() {
			So(f.Validate(cfg), ShouldBeNil)
		})
	})

	Convey("Given a Feedback model where 'ons_url' is subnet of domain", t, func() {
		f := validFeedbackModel()
		f.OnsURL = fmt.Sprintf("https://www.%s:1234/sub/path", onsHost)

		Convey("Then validation is successful", func() {
			So(f.Validate(cfg), ShouldBeNil)
		})
	})

	Convey("Given a Feedback model where 'ons_url' is subnet of invalid domain", t, func() {
		f := validFeedbackModel()
		f.OnsURL = "https://www.somedomain/sub/path"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Key: 'Feedback.OnsURL' Error:Field validation for 'OnsURL' failed on the 'ons_url' tag")
		})
	})

	Convey("Given a Feedback model where 'is_page_useful' is not provided'", t, func() {
		f := validFeedbackModel()
		f.IsPageUseful = nil

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Key: 'Feedback.IsPageUseful' Error:Field validation for 'IsPageUseful' failed on the 'required' tag")
		})
	})

	Convey("Given a Feedback model where 'is_general_feedback' is not provided'", t, func() {
		f := validFeedbackModel()
		f.IsGeneralFeedback = nil

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Key: 'Feedback.IsGeneralFeedback' Error:Field validation for 'IsGeneralFeedback' failed on the 'required' tag")
		})
	})

	Convey("Given a Feedback model where 'email_address' has the wrong format'", t, func() {
		f := validFeedbackModel()
		f.EmailAddress = "thisIsNotAnEmail"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Key: 'Feedback.EmailAddress' Error:Field validation for 'EmailAddress' failed on the 'email' tag")
		})
	})

	Convey("Given a Feedback model with an invalid 'ons_url' value", t, func() {
		f := validFeedbackModel()
		f.OnsURL = "£@%"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Key: 'Feedback.OnsURL' Error:Field validation for 'OnsURL' failed on the 'ons_url' tag")
		})
	})

	Convey("Given a Feedback model with an unexpected 'ons_url' value", t, func() {
		f := validFeedbackModel()
		f.OnsURL = "http://attackerHost:1234/some/path"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Key: 'Feedback.OnsURL' Error:Field validation for 'OnsURL' failed on the 'ons_url' tag")
		})
	})
}

func TestFeedbackSanitize(t *testing.T) {
	Convey("Given a Feedback model where all strings are unsafe", t, func() {
		f := validFeedbackModel()
		f.OnsURL = unsafeStr
		f.Feedback = unsafeStr
		f.Name = unsafeStr
		f.EmailAddress = unsafeStr

		Convey("Then sanitize with the config enabled mutates the model accordingly", func() {
			cfg := &config.Sanitize{
				HTML:  true,
				SQL:   true,
				NoSQL: true,
			}

			expected := validFeedbackModel()
			expected.OnsURL = sanitizedStr
			expected.Feedback = sanitizedStr
			expected.Name = sanitizedStr
			expected.EmailAddress = sanitizedStr

			f.Sanitize(cfg)
			So(f, ShouldResemble, expected)
		})

		Convey("Then sanitize with the config disabled does not mutate the model", func() {
			cfg := &config.Sanitize{
				HTML:  false,
				SQL:   false,
				NoSQL: false,
			}

			expected := validFeedbackModel()
			expected.OnsURL = unsafeStr
			expected.Feedback = unsafeStr
			expected.Name = unsafeStr
			expected.EmailAddress = unsafeStr

			f.Sanitize(cfg)
			So(f, ShouldResemble, expected)
		})
	})
}
