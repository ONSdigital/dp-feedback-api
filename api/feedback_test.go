package api_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ONSdigital/dp-feedback-api/api"
	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/gorilla/mux"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	ctx             = context.Background()
	cfg             = &config.Config{OnsDomain: "testhost"}
	feedbackPayload = `{
		"is_page_useful": true,
		"is_general_feedback": false,
		"ons_url": "https://testhost:1234/sub/path",
		"feedback": "very nice and useful website!",
		"name": "Mr Feedback reporter",
		"email_address": "feedback@reporter.com"
	}`
)

func TestPostFeedbackHandler(t *testing.T) {
	Convey("Given an API", t, func() {
		a := api.Setup(ctx, cfg, mux.NewRouter(), nil)

		Convey("A post feedback request with a valid body results in 200 OK response with empty response body", func() {
			b := bytes.NewBufferString(feedbackPayload)
			req := httptest.NewRequest(http.MethodPost, "http://localhost:8080/feedback", b)
			resp := httptest.NewRecorder()

			a.Router.ServeHTTP(resp, req)

			So(resp.Code, ShouldEqual, http.StatusOK)
			So(resp.Body.Len(), ShouldEqual, 0)
		})
	})
}
