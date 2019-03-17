package config

import (
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

// App represent the global application configuration.
type App struct {
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

// Uptime returns the uptime duration since the time it was deployed.
func (a *App) Uptime() string {
	return time.Since(a.StartAt).String()
}

// IsProduction returns true if the current environment is set to "production".
func (a *App) IsProduction() bool {
	return a.Env == "production"
}

// IsDevelopment returns true if the current environment is set to
// "development".
func (a *App) IsDevelopment() bool {
	return a.Env == "development"
}

// New returns a new Config pointer with populated values. Will panic if the
// required environment variables are not set.
func New() *App {
	var app App
	err := envconfig.Process("", &app)
	if err != nil {
		log.Fatal(err)
	}
	app.StartAt = time.Now()

	hostname, err := os.Hostname()
	if err != nil {
		log.Fatal(err)
	}
	app.Hostname = hostname
	return &app
}
