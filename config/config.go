package config

import (
	"log"
	"os"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BuildDate time.Time `envconfig:"BUILD_DATE"`
	Env       string    `envconfig:"ENV" default:"development"`
	Hostname  string    `ignored:"true"`
	Port      string    `envconfig:"PORT" default:"8080"`
	StartAt   time.Time `ignored:"true"`
	Tag       string    `envconfig:"TAG"`
}

func (c *Config) Uptime() string {
	return time.Since(c.StartAt).String()
}

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
