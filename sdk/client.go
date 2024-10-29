package sdk

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	healthcheck "github.com/ONSdigital/dp-api-clients-go/v2/health"
	"github.com/ONSdigital/dp-feedback-api/models"
	sdkError "github.com/ONSdigital/dp-feedback-api/sdk/errors"
	health "github.com/ONSdigital/dp-healthcheck/healthcheck"
)

// package level constants
const (
	Service          = "dp-feedback-api"
	FeedbackEndpoint = "%s/v1/feedback"
	Authorization    = "Authorization"
	BearerPrefix     = "Bearer "
)

// HTTPClient is the interface that defines a client for making HTTP requests
type HTTPClient interface {
	Do(ctx context.Context, req *http.Request) (*http.Response, error)
}

// APIClient implementation of permissions.Store that gets permission data from the permissions API
type Client struct {
	hcCli *healthcheck.Client
}

// Options is a struct containing for customised options for the API client
type Options struct {
	AuthToken string
}

func (o *Options) SetAuth(req *http.Request) {
	if o.AuthToken != "" {
		req.Header.Add(Authorization, fmt.Sprintf("%s%s", BearerPrefix, o.AuthToken))
	}
}

// New constructs a new Client instance with a given feedback api url
func New(feedbackAPIURL string) *Client {
	return &Client{
		hcCli: healthcheck.NewClient(Service, feedbackAPIURL),
	}
}

// NewWithHealthClient creates a new instance of search API Client,
// reusing the URL and Clienter from the provided healthcheck client
func NewWithHealthClient(hcCli *healthcheck.Client) *Client {
	return &Client{
		hcCli: healthcheck.NewClientWithClienter(Service, hcCli.URL, hcCli.Client),
	}
}

// URL returns the URL used by this client
func (cli *Client) URL() string {
	return cli.hcCli.URL
}

// Health returns the underlying Healthcheck Client for this search API client
func (cli *Client) Health() *healthcheck.Client {
	return cli.hcCli
}

// Checker calls search api health endpoint and returns a check object to the caller
func (cli *Client) Checker(ctx context.Context, check *health.CheckState) error {
	return cli.hcCli.Checker(ctx, check)
}

// PostFeedback sends the provided feedback model to the feedback API via a post call
func (cli *Client) PostFeedback(ctx context.Context, feedback *models.Feedback, options Options) *sdkError.StatusError {
	uri := fmt.Sprintf(FeedbackEndpoint, cli.hcCli.URL)

	buf := &bytes.Buffer{}
	err := json.NewEncoder(buf).Encode(feedback)
	if err != nil {
		return &sdkError.StatusError{
			Err:  fmt.Errorf("failed to encode feedback: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	req, err := http.NewRequest(http.MethodPost, uri, buf)
	if err != nil {
		return &sdkError.StatusError{
			Err:  fmt.Errorf("error creating request: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	options.SetAuth(req)

	resp, err := cli.hcCli.Client.Do(ctx, req)
	if err != nil {
		return &sdkError.StatusError{
			Err:  fmt.Errorf("error sending request: %w", err),
			Code: http.StatusInternalServerError,
		}
	}

	defer func() {
		if resp.Body != nil {
			resp.Body.Close()
		}
	}()

	if resp.StatusCode != http.StatusCreated {
		return &sdkError.StatusError{
			Err:  fmt.Errorf("unexpected status returned from the feedback api post feedback endpoint: %d", resp.StatusCode),
			Code: resp.StatusCode,
		}
	}

	return nil
}
