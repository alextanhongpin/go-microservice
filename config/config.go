package config

import (
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represent the global application configuration.
type Config struct {
	BuildDate time.Time `envconfig:"BUILD_DATE"`
	Env       string    `envconfig:"ENV" default:"development"`
	Port      string    `envconfig:"PORT" default:"8080"`
	Tag       string    `envconfig:"TAG"`
	Hostname  string    `ignored:"true"`
	StartAt   time.Time `ignored:"true"`
}

// Uptime returns the uptime duration since the time it was deployed.
func (c *Config) Uptime() string {
	return time.Since(c.StartAt).String()
}

// IsProduction returns true if the current environment is set to "production".
func (c *Config) IsProduction() bool {
	return c.Env == "production"
}

// IsDevelopment returns true if the current environment is set to
// "development".
func (c *Config) IsDevelopment() bool {
	return c.Env == "development"
}

// New returns a new Config pointer with populated values. Will panic if the
// required environment variables are not set.
func New() *Config {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		log.Fatal(err)
	}
	cfg.StartAt = time.Now()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	cfg.Hostname = hostname
	return &cfg
}
