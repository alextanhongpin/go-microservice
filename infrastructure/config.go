package infrastructure

import (
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// Config represent the global application configuration.
type Config struct {
	Name       string    `envconfig:"NAME" default:"yourapp"`
	Audience   string    `envconfig:"AUDIENCE" required:"true"`
	BuildDate  time.Time `envconfig:"BUILD_DATE"`
	Credential string    `envconfig:"CREDENTIAL" required:"true"`
	Env        string    `envconfig:"ENV" default:"development"`
	Issuer     string    `envconfig:"ISSUER" required:"true"`
	Port       string    `envconfig:"PORT" default:"8080"`
	Secret     string    `envconfig:"SECRET" required:"true"`
	Semver     string    `envconfig:"SEMVER" required:"true"`
	Tag        string    `envconfig:"TAG"`
	Hostname   string    `ignored:"true"`
	StartAt    time.Time `ignored:"true"`

	// Nested Option.
	Database
}

type Database struct {
	User string `envconfig:"DB_USER" required:"true"`
	Pass string `envconfig:"DB_PASS" required:"true"`
	Host string `envconfig:"DB_HOST" required:"true"`
	Name string `envconfig:"DB_NAME" required:"true"`
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

// NewConfig returns a new Config pointer with populated values. Will panic if the
// required environment variables are not set.
func NewConfig() *Config {
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
