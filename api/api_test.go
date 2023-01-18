package api

import (
	"context"
	"net/http"
	"testing"

	"github.com/ONSdigital/dp-feedback-api/config"
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
		api := Setup(ctx, cfg, r, nil, nil)

		Convey("When created the following routes should have been added", func() {
			So(hasRoute(api.Router, "/feedback", http.MethodPost), ShouldBeTrue)
		})
	})
}

func hasRoute(r chi.Router, path, method string) bool {
	return r.Match(chi.NewRouteContext(), method, path)
}
