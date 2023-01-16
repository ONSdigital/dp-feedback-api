package config

import (
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represents service configuration for dp-feedback-api
type Config struct {
	BindAddr                   string        `envconfig:"BIND_ADDR"`
	GracefulShutdownTimeout    time.Duration `envconfig:"GRACEFUL_SHUTDOWN_TIMEOUT"`
	HealthCheckInterval        time.Duration `envconfig:"HEALTHCHECK_INTERVAL"`
	HealthCheckCriticalTimeout time.Duration `envconfig:"HEALTHCHECK_CRITICAL_TIMEOUT"`
	OnsDomain                  string        `envconfig:"ONS_DOMAIN"`
	Mail                       *Mail
}

// MailConfig represents the subset of configuration corresponding to the email service
type Mail struct {
	Host     string `envconfig:"MAIL_HOST"`
	User     string `envconfig:"MAIL_USER"`
	Password string `envconfig:"MAIL_PASSWORD" json:"-"`
	Port     string `envconfig:"MAIL_PORT"`
}

var cfg *Config

// Get returns the default config with any modifications through environment
// variables
func Get() (*Config, error) {
	if cfg != nil {
		return cfg, nil
	}

	cfg = &Config{
		BindAddr:                   "localhost:28600",
		GracefulShutdownTimeout:    5 * time.Second,
		HealthCheckInterval:        30 * time.Second,
		HealthCheckCriticalTimeout: 90 * time.Second,
		OnsDomain:                  "localhost",
		Mail: &Mail{
			Host:     "localhost",
			Port:     "1025",
			User:     "",
			Password: "",
		},
	}

	return cfg, envconfig.Process("", cfg)
}

func (c *Config) Validate() error {
	return nil
}
