package models_test

import (
	"context"
	"errors"
	"fmt"
	"testing"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx               = context.Background()
	onsHost           = "testhost"
	cfg               = &config.Config{OnsDomain: onsHost}
	pageIsUseful      = true
	isGeneralFeedback = true
	unsafeHTML        = `<script>document.getElementById("demo").innerHTML = "Hello JavaScript!";</script>`
	sanitizedHTML     = `&lt;script&gt;document.getElementById(&#34;demo&#34;).innerHTML = &#34;Hello JavaScript!&#34;;&lt;/script&gt;`
	unsafeSQL         = `TODO`
	sanitizedSQL      = `TODO`
	unsafeNoSQL       = `TODO`
	sanitizedNoSQL    = `TODO`
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

		Convey("Then the validation is successful", func() {
			So(f.Validate(cfg), ShouldBeNil)
		})
	})

	Convey("Given a valid Feedback model containing only the mandatory fields", t, func() {
		f := &models.Feedback{
			IsPageUseful:      &pageIsUseful,
			IsGeneralFeedback: &isGeneralFeedback,
		}

		Convey("Then the validation is successful", func() {
			So(f.Validate(cfg), ShouldBeNil)
		})
	})

	Convey("Given a Feedback model where 'is_page_useful' is not provided'", t, func() {
		f := validFeedbackModel()
		f.IsPageUseful = nil

		Convey("Then the validation fails with the expected error", func() {
			So(f.Validate(cfg), ShouldResemble, errors.New("is_page_useful is compulsory"))
		})
	})

	Convey("Given a Feedback model where 'is_general_feedback' is not provided'", t, func() {
		f := validFeedbackModel()
		f.IsGeneralFeedback = nil

		Convey("Then the validation fails with the expected error", func() {
			So(f.Validate(cfg), ShouldResemble, errors.New("is_general_feedback is compulsory"))
		})
	})

	Convey("Given a Feedback model where 'email_address' has the wrong format'", t, func() {
		f := validFeedbackModel()
		f.EmailAddress = "thisIsNotAnEmail"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "invalid email address: mail: missing '@' or angle-addr")
		})
	})

	Convey("Given a Feedback model with an invalid 'ons_url' value", t, func() {
		f := validFeedbackModel()
		f.OnsURL = "£@%"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, `invalid ons url: parse "£@%": invalid URL escape "%"`)
		})
	})

	Convey("Given a Feedback model with an unexpected 'ons_url' value", t, func() {
		f := validFeedbackModel()
		f.OnsURL = "http://attackerHost:1234/some/path"

		Convey("Then the validation fails with the expected error", func() {
			err := f.Validate(cfg)
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "unexpected ons domain name: attackerHost")
		})
	})
}

func TestFeedbackSanitize(t *testing.T) {

	Convey("Given a Feedback model where all strings are unsafe", t, func() {

		f := validFeedbackModel()
		f.OnsURL = unsafeHTML
		f.Feedback = unsafeHTML
		f.Name = unsafeHTML
		f.EmailAddress = unsafeHTML

		Convey("Then sanitize mutates the model accordingly", func() {
			expected := validFeedbackModel()
			expected.OnsURL = sanitizedHTML
			expected.Feedback = sanitizedHTML
			expected.Name = sanitizedHTML
			expected.EmailAddress = sanitizedHTML

			f.Sanitize()
			So(f, ShouldResemble, expected)
		})
	})
}