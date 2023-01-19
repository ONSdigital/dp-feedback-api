package api_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-feedback-api/api"
	"github.com/ONSdigital/dp-feedback-api/api/mock"
	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/models"
	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx = context.Background()
	cfg = &config.Config{
		OnsDomain: "testhost",
		Mail: &config.Mail{
			From: "sender@mail.com",
			To:   "receiver@mail.com",
		},
	}
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

Is page useful: true
Is general feedback: false
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

// TestPostFeedbackHandler - this tests the full handler
// TODO this may be removed according to the task requirements, and implemented as component tasks
func TestPostFeedbackHandler(t *testing.T) {

	var sendOK = func(from string, to []string, msg []byte) error {
		return nil
	}

	var sendFail = func(from string, to []string, msg []byte) error {
		return errors.New("failed to send email")
	}

	Convey("Given an API with the PostFeedback handler registered", t, func() {
		a := &api.API{
			Cfg:    cfg,
			Router: chi.NewRouter(),
		}
		a.Router.Post("/feedback", a.PostFeedback)

		Convey("And a successful email sender mock", func() {
			emailSenderMock := &mock.EmailSenderMock{SendFunc: sendOK}
			a.EmailSender = emailSenderMock

			Convey("When a post feedback request with a valid body is sent", func() {
				b := bytes.NewBufferString(feedbackPayload)
				req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/feedback", b)
				resp := httptest.NewRecorder()

				a.Router.ServeHTTP(resp, req)

				Convey("Then a response with status 201 CREATED and empty body is received", func() {
					So(resp.Code, ShouldEqual, http.StatusCreated)
					So(resp.Body.Len(), ShouldEqual, 0)
				})

				Convey("Then the expected email is sent", func() {
					So(emailSenderMock.SendCalls(), ShouldHaveLength, 1)
					So(emailSenderMock.SendCalls()[0].From, ShouldEqual, cfg.Mail.From)
					So(emailSenderMock.SendCalls()[0].To, ShouldResemble, []string{cfg.Mail.To})
					So(string(emailSenderMock.SendCalls()[0].Msg), ShouldEqual, expectedEmail)
				})
			})
		})

		Convey("And an unsuccessful email sender mock", func() {
			emailSenderMock := &mock.EmailSenderMock{SendFunc: sendFail}
			a.EmailSender = emailSenderMock

			Convey("When a post feedback request with a valid body is sent", func() {
				b := bytes.NewBufferString(feedbackPayload)
				req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/feedback", b)
				resp := httptest.NewRecorder()

				a.Router.ServeHTTP(resp, req)

				Convey("Then a response with status 500 INTERNAL SERVER ERROR is received", func() {
					So(resp.Code, ShouldEqual, http.StatusInternalServerError)
				})
			})
		})

		Convey("When a post feedback request with a wrong body is sent", func() {
			malformed := `{"this" is not a valid json`
			b := bytes.NewBufferString(malformed)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/feedback", b)
			resp := httptest.NewRecorder()

			a.Router.ServeHTTP(resp, req)

			Convey("Then a response with status 400 BAD REQUEST is received", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})

		Convey("When a post feedback request with an invalid feedback body is sent", func() {
			invalid := `{
				"is_page_useful": false,
				"name": "Mr Feedback reporter"
			}`
			b := bytes.NewBufferString(invalid)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/feedback", b)
			resp := httptest.NewRecorder()

			a.Router.ServeHTTP(resp, req)

			Convey("Then a response with status 400 BAD REQUEST is received", func() {
				So(resp.Code, ShouldEqual, http.StatusBadRequest)
			})
		})
	})
}

func TestGenerateFeedbackMessage(t *testing.T) {
	Convey("The expected email is generated from a valid feedback model", t, func() {
		f := testFeedback()
		generated := api.GenerateFeedbackMessage(f, "sender@mail.com", "receiver@mail.com")
		So(string(generated), ShouldEqual, expectedEmail)
	})
}
