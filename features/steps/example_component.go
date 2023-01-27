package steps

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	componenttest "github.com/ONSdigital/dp-component-test"
	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/service"
	"github.com/ONSdigital/dp-feedback-api/service/mock"
)

var (
	BuildTime = strconv.Itoa(time.Now().Nanosecond())
	GitCommit = "component test commit"
	Version   = "component test version"
)

type Component struct {
	componenttest.ErrorFeature
	svc            *service.Service
	errorChan      chan error
	Config         *config.Config
	HTTPServer     *http.Server
	ServiceRunning bool
	apiFeature     *componenttest.APIFeature
}

func NewComponent() (*Component, error) {
	c := &Component{
		HTTPServer:     &http.Server{},
		errorChan:      make(chan error),
		ServiceRunning: false,
	}

	var err error
	c.Config, err = config.Get()
	if err != nil {
		return nil, err
	}

	c.apiFeature = componenttest.NewAPIFeature(c.Router)
	c.setInitialiserMock()
	c.svc = service.New()
	c.svc.Config = c.Config

	return c, nil
}

// Router initialises the service, returning the service's (server) router for tests
// This delayed initialisation is needed to ensure that any changes to the router (or the service in general)
// as a result of test setup, are picked up
func (c *Component) Router() (http.Handler, error) {
	if err := c.svc.Init(context.Background(), c.svc.Config, BuildTime, GitCommit, Version); err != nil {
		return nil, fmt.Errorf("failed to initialise service: %s", err)
	}

	return c.svc.API.Router, nil
}

func (c *Component) Reset() *Component {
	c.apiFeature.Reset()
	return c
}

func (c *Component) Close() error {
	if c.svc != nil && c.ServiceRunning {
		c.svc.Close(context.Background())
		c.ServiceRunning = false
	}
	return nil
}

func (c *Component) setInitialiserMock() {
	service.GetHTTPServer = func(bindAddr string, router http.Handler) service.HTTPServer {
		return &http.Server{Addr: bindAddr, Handler: router}
	}

	service.GetEmailSender = func(*config.Mail) service.EmailSender {
		return &mock.EmailSenderMock{}
	}
}

// func (c *Component) InitialiseService() (http.Handler, error) {
// 	var err error
// 	c.svc, err = service.Run(context.Background(), c.Config, c.svcList, "1", "", "", c.errorChan)
// 	if err != nil {
// 		return nil, err
// 	}

// 	c.ServiceRunning = true
// 	return c.HTTPServer.Handler, nil
// }

// func (c *Component) DoGetHealthcheckOk(cfg *config.Config, buildTime, gitCommit, version string) (service.HealthChecker, error) {
// 	return &mock.HealthCheckerMock{
// 		AddCheckFunc: func(name string, checker healthcheck.Checker) error { return nil },
// 		StartFunc:    func(ctx context.Context) {},
// 		StopFunc:     func() {},
// 	}, nil
// }

// func (c *Component) DoGetHTTPServer(bindAddr string, router http.Handler) service.HTTPServer {
// 	c.HTTPServer.Addr = bindAddr
// 	c.HTTPServer.Handler = router
// 	return c.HTTPServer
// }
