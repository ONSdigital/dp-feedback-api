package api_test

import (
	"testing"

	"github.com/ONSdigital/dp-feedback-api/api"
	"github.com/ONSdigital/dp-feedback-api/models"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	feedbackPayload = `{
		"is_page_useful": true,
		"is_general_feedback": false,
		"ons_url": "https://testhost:1234/sub/path",
		"feedback": "very nice and useful website!",
		"name": "Mr Feedback reporter",
		"email_address": "feedback@reporter.com"
	}`
	isPageUseful      = true
	isGeneralFeedback = false
)

var expectedEmail = `From: sender@mail.com
To: receiver@mail.com
Subject: Feedback received

Feedback Type: A specific page
Page URL: https://testhost:1234/sub/path
Description: very nice and useful website!
Name: Mr Feedback reporter
Email address: feedback@reporter.com
`

func testFeedback() *models.Feedback {
	return &models.Feedback{
		IsPageUseful:      &isPageUseful,
		IsGeneralFeedback: &isGeneralFeedback,
		OnsURL:            "https://testhost:1234/sub/path",
		Feedback:          "very nice and useful website!",
		Name:              "Mr Feedback reporter",
		EmailAddress:      "feedback@reporter.com",
	}
}

func TestGenerateFeedbackMessage(t *testing.T) {
	Convey("The expected email is generated from a valid feedback model", t, func() {
		f := testFeedback()
		generated := api.GenerateFeedbackMessage(f, "sender@mail.com", "receiver@mail.com")
		So(string(generated), ShouldEqual, expectedEmail)
	})
}
