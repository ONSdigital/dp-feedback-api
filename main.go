package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ONSdigital/dp-feedback-api/config"
	"github.com/ONSdigital/dp-feedback-api/service"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/pkg/errors"
)

const serviceName = "dp-feedback-api"

var (
	// BuildTime represents the time in which the service was built
	BuildTime string
	// GitCommit represents the commit (SHA-1) hash of the service that is running
	GitCommit string
	// Version represents the version of the service that is running
	Version string

// TODO: remove below explainer before commiting
/* NOTE: replace the above with the below to run code with for example vscode debugger.
BuildTime string = "1601119818"
GitCommit string = "6584b786caac36b6214ffe04bf62f058d4021538"
Version   string = "v0.1.0"
*/
)

func main() {
	log.Namespace = serviceName
	ctx := context.Background()

	if err := run(ctx); err != nil {
		log.Fatal(ctx, "fatal runtime error", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	svcErrors := make(chan error, 1)

	// Read config
	cfg, err := config.Get()
	if err != nil {
		return errors.Wrap(err, "error getting configuration")
	}
	log.Info(ctx, "config on startup", log.Data{"config": cfg, "build_time": BuildTime, "git-commit": GitCommit})

	// Make sure that context is cancelled when 'run' finishes its execution.
	// Any remaining go-routine that was not terminated during svc.Close (graceful shutdown) will be terminated by ctx.Done()
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	// Run the service
	svc := service.New()
	if err := svc.Init(ctx, cfg, BuildTime, GitCommit, Version); err != nil {
		return errors.Wrap(err, "running service failed")
	}
	svc.Start(ctx, svcErrors)

	// blocks until an os interrupt or a fatal error occurs
	select {
	case err := <-svcErrors:
		err = fmt.Errorf("service error received: %w", err)
		if errClose := svc.Close(ctx); errClose != nil {
			log.Error(ctx, "service close error during error handling", errClose)
		}
		return errors.Wrap(err, "service error received")
	case sig := <-signals:
		log.Info(ctx, "os signal received", log.Data{"signal": sig})
	}
	return svc.Close(ctx)
}
