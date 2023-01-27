package service

import (
	"context"
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-feedback-api/api"
	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

// Service contains all the configs, server and clients to run the API
type Service struct {
	Config      *config.Config
	Server      HTTPServer
	API         *api.API
	EmailSender EmailSender
	HealthCheck HealthChecker
}

func New() *Service {
	return &Service{}
}

func (svc *Service) Init(ctx context.Context, cfg *config.Config, buildTime, gitCommit, version string) error {
	log.Info(ctx, "initialising service")
	var err error

	// Check config
	if cfg == nil {
		return errors.New("nil config passed to service init")
	}
	svc.Config = cfg

	// Get Email Sender
	svc.EmailSender = GetEmailSender(cfg.Mail)

	// Get HealthCheck
	if svc.HealthCheck, err = GetHealthCheck(cfg, buildTime, gitCommit, version); err != nil {
		return fmt.Errorf("could not instantiate healthcheck: %w", err)
	}
	if err := svc.registerCheckers(ctx); err != nil {
		return fmt.Errorf("unable to register checkers: %w", err)
	}

	// Create an HTTP server containing a new router with /health endpoint
	r := chi.NewRouter()
	r.Handle("/health", http.HandlerFunc(svc.HealthCheck.Handler))
	svc.Server = GetHTTPServer(cfg.BindAddr, r)

	// Create API
	svc.API = api.Setup(ctx, cfg, r, svc.EmailSender)
	return nil
}

// Start the service
func (svc *Service) Start(ctx context.Context, svcErrors chan error) {
	log.Info(ctx, "starting service")

	svc.HealthCheck.Start(ctx)

	// Run the http server in a new go-routine
	go func() {
		if err := svc.Server.ListenAndServe(); err != nil {
			svcErrors <- errors.Wrap(err, "failure in http listen and serve")
		}
	}()
}

// Close gracefully shuts the service down in the required order, with timeout
func (svc *Service) Close(ctx context.Context) error {
	timeout := svc.Config.GracefulShutdownTimeout
	log.Info(ctx, "commencing graceful shutdown", log.Data{"graceful_shutdown_timeout": timeout})
	ctx, cancel := context.WithTimeout(ctx, timeout)

	// track shutown gracefully closes up
	var hasShutdownError bool

	go func() {
		defer cancel()

		// stop healthcheck, as it depends on everything else
		if svc.HealthCheck != nil {
			log.Info(ctx, "stopping health checker...")
			svc.HealthCheck.Stop()
			log.Info(ctx, "successfully stopped health checker")
		}

		// stop any incoming requests before closing any outbound connections
		if svc.Server != nil {
			log.Info(ctx, "stopping http server...")
			if err := svc.Server.Shutdown(ctx); err != nil {
				log.Error(ctx, "failed to shutdown http server", err)
				hasShutdownError = true
			}
			log.Info(ctx, "successfully stopped http server")
		}

		// TODO: Close other dependencies, in the expected order
	}()

	// wait for shutdown success (via cancel) or failure (timeout)
	<-ctx.Done()

	// timeout expired
	if ctx.Err() == context.DeadlineExceeded {
		log.Error(ctx, "shutdown timed out", ctx.Err())
		return ctx.Err()
	}

	// other error
	if hasShutdownError {
		return fmt.Errorf("failed to shutdown gracefully")
	}

	log.Info(ctx, "graceful shutdown was successful")
	return nil
}

func (svc *Service) registerCheckers(ctx context.Context) (err error) {
	if svc.HealthCheck == nil {
		return errors.New("healthcheck must be created before registering checkers")
	}

	return nil
}
