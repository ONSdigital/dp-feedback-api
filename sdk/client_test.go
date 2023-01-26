package sdk_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
	"time"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-feedback-api/models"
	"github.com/ONSdigital/dp-feedback-api/sdk"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
	dphttp "github.com/ONSdigital/dp-net/v2/http"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testHost = "http://localhost:1234"
)

var (
	initialTestState = healthcheck.CreateCheckState(sdk.Service)
	testAuthToken    = "serviceToken"
)

func TestNewClient(t *testing.T) {
	Convey("Given some host", t, func() {
		someHost := testHost
		Convey("When NewClient is called", func() {
			apiClient := sdk.New(someHost)
			Convey("Then the api is not nil", func() {
				So(apiClient, ShouldNotBeNil)
			})
		})
	})
}

func TestNewClientWithOptions(t *testing.T) {
	Convey("Given some host and custom options", t, func() {
		hcCli, _ := getMockClient(testHost, http.StatusCreated, "", nil)
		Convey("When NewClientWithOptions is called", func() {
			apiClient := sdk.NewWithHealthClient(hcCli)
			Convey("Then the api is not nil", func() {
				So(apiClient, ShouldNotBeNil)
			})
		})
	})
}

func TestURL(t *testing.T) {
	Convey("Given a mock http Client", t, func() {
		hcCli, _ := getMockClient(testHost, http.StatusCreated, "", nil)
		apiClient := sdk.NewWithHealthClient(hcCli)

		Convey("Then calling URL() returns the expected host", func() {
			So(apiClient.URL(), ShouldEqual, testHost)
		})
	})
}

func TestHealth(t *testing.T) {
	Convey("Given a mock http Client", t, func() {
		hcCli, _ := getMockClient(testHost, http.StatusCreated, "", nil)
		apiClient := sdk.NewWithHealthClient(hcCli)

		Convey("Then calling Health() returns the expected Health client", func() {
			So(apiClient.Health(), ShouldResemble, hcCli)
		})
	})
}

func TestHealthChecker(t *testing.T) {
	ctx := context.Background()
	timePriorHealthCheck := time.Now().UTC()
	path := "/health"

	Convey("Given clienter.Do returns an error", t, func() {
		doErr := errors.New("unexpected error")
		hcCli, httpClientMock := getMockClient(testHost, -1, "", doErr)
		apiClient := sdk.NewWithHealthClient(hcCli)
		check := initialTestState

		Convey("When feedback API client Checker is called", func() {
			err := apiClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, sdk.Service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 0)
				So(check.Message(), ShouldEqual, doErr.Error())
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClientMock.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})

	Convey("Given a 500 response for health check", t, func() {
		hcCli, httpClientMock := getMockClient(testHost, http.StatusInternalServerError, "internal error", nil)
		apiClient := sdk.NewWithHealthClient(hcCli)
		check := initialTestState

		Convey("When feedback API client Checker is called", func() {
			err := apiClient.Checker(ctx, &check)
			So(err, ShouldBeNil)

			Convey("Then the expected check is returned", func() {
				So(check.Name(), ShouldEqual, sdk.Service)
				So(check.Status(), ShouldEqual, health.StatusCritical)
				So(check.StatusCode(), ShouldEqual, 500)
				So(check.Message(), ShouldEqual, sdk.Service+healthcheck.StatusMessage[health.StatusCritical])
				So(*check.LastChecked(), ShouldHappenAfter, timePriorHealthCheck)
				So(check.LastSuccess(), ShouldBeNil)
				So(*check.LastFailure(), ShouldHappenAfter, timePriorHealthCheck)
			})

			Convey("And client.Do should be called once with the expected parameters", func() {
				doCalls := httpClientMock.DoCalls()
				So(doCalls, ShouldHaveLength, 1)
				So(doCalls[0].Req.URL.Path, ShouldEqual, path)
			})
		})
	})
}

func TestPostFeedback(t *testing.T) {
	Convey("Given a mock http client that returns 201 created", t, func() {
		hcCli, httpClientMock := getMockClient(testHost, http.StatusCreated, "", nil)
		apiClient := sdk.NewWithHealthClient(hcCli)

		Convey("When PostFeedback is called with a valid feedback body", func() {
			ctx := context.Background()
			f := getExampleFeedback()
			opts := sdk.Options{AuthToken: testAuthToken}
			err := apiClient.PostFeedback(ctx, f, opts)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("Then the expected request is sent with the expected path, method and auth header", func() {
				So(httpClientMock.DoCalls(), ShouldHaveLength, 1)
				So(httpClientMock.DoCalls()[0].Req.URL.String(), ShouldEqual, "http://localhost:1234/v1/feedback")
				So(httpClientMock.DoCalls()[0].Req.Method, ShouldEqual, http.MethodPost)
				So(httpClientMock.DoCalls()[0].Req.Header.Get(sdk.Authorization), ShouldEqual, "Bearer serviceToken")
			})
		})
	})

	Convey("Given a mock http client that returns 401 Unauthorized", t, func() {
		hcCli, _ := getMockClient(testHost, http.StatusUnauthorized, "401 Unauthorized", nil)
		apiClient := sdk.NewWithHealthClient(hcCli)

		Convey("When PostFeedback is called with a valid feedback body", func() {
			ctx := context.Background()
			f := getExampleFeedback()
			opts := sdk.Options{AuthToken: "wrong"}
			err := apiClient.PostFeedback(ctx, f, opts)

			Convey("Then the expected error and status code is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "unexpected status returned from the feedback api post feedback endpoint: 401")
				So(err.Status(), ShouldEqual, http.StatusUnauthorized)
			})
		})
	})

	Convey("Given a mock http client that fails with an unexpected error", t, func() {
		doErr := errors.New("unexpected error")
		hcCli, _ := getMockClient(testHost, -1, "", doErr)
		apiClient := sdk.NewWithHealthClient(hcCli)

		Convey("When PostFeedback is called with a valid feedback body", func() {
			ctx := context.Background()
			f := getExampleFeedback()
			opts := sdk.Options{AuthToken: testAuthToken}
			err := apiClient.PostFeedback(ctx, f, opts)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "error sending request: unexpected error")
			})
		})
	})
}

func getMockClient(host string, statusCode int, bodyStr string, doErr error) (*healthcheck.Client, *dphttp.ClienterMock) {
	c := &dphttp.ClienterMock{
		DoFunc: func(ctx context.Context, req *http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: statusCode,
				Body:       body(bodyStr),
			}, doErr
		},
		GetPathsWithNoRetriesFunc: func() []string {
			return []string{"/healthcheck"}
		},
		SetPathsWithNoRetriesFunc: func(paths []string) {},
	}
	return healthcheck.NewClientWithClienter(sdk.Service, host, c), c
}

func getExampleFeedbackJson() []byte {
	f := getExampleFeedback()
	feedbackJson, err := json.Marshal(f)
	So(err, ShouldBeNil)
	return feedbackJson
}

func getExampleFeedback() *models.Feedback {
	pageUseful := true
	generalFeedback := true

	return &models.Feedback{
		IsPageUseful:      &pageUseful,
		IsGeneralFeedback: &generalFeedback,
	}
}

func body(strBody string) io.ReadCloser {
	buff := bytes.NewBufferString(strBody)
	return io.NopCloser(buff)
}
