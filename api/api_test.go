package api_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-feedback-api/api"
	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/models"
	"github.com/go-chi/chi/v5"
	. "github.com/smartystreets/goconvey/convey"
)

func TestSetup(t *testing.T) {
	Convey("Given an API instance", t, func() {
		r := chi.NewRouter()
		ctx := context.Background()
		cfg := &config.Config{
			OnsDomain: "localhost",
		}
		a := api.Setup(ctx, cfg, r, nil, nil)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(a.Router, "/feedback", http.MethodPost), ShouldBeTrue)
		})
	})
}

func TestUnmarshal(t *testing.T) {
	Convey("A valid feedback payload body is correctly unmarshaled to a Feedback model", t, func() {
		b := body(feedbackPayload)
		target := &models.Feedback{}
		expected := testFeedback()
		err := api.Unmarshal(b, target)
		So(err, ShouldBeNil)
		So(target, ShouldResemble, expected)
	})

	Convey("An incorrect feedback payload body fails to unmarshal to a Feedback model", t, func() {
		b := body("{'this' is not a valid json")
		target := &models.Feedback{}
		err := api.Unmarshal(b, target)
		So(err, ShouldNotBeNil)
	})
}

func body(strBody string) io.ReadCloser {
	buff := bytes.NewBufferString(strBody)
	return io.NopCloser(buff)
}

func hasRoute(r chi.Router, path, method string) bool {
	return r.Match(chi.NewRouteContext(), method, path)
}
